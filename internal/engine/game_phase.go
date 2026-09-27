package engine

// This file holds the phase machine (ADR 0012): the loop that walks a turn
// through its eight phases and the body of each engine-driven phase. The turn
// lifecycle entry points that resume the loop after a player decision —
// StartTurn, ChooseHouse, EndPlayPhase — live in game_turn.go.

// Phase is the part of the turn now running, so a frontend can draw the turn's
// skeleton and a card can ask which phase it is resolving in.
func (g *Game) Phase() Phase { return g.State.Phase }

// runPhases advances the turn from the current phase, running each engine-driven
// phase to completion, and stops when it reaches a phase that waits for the
// frontend, when the last phase of the turn is done, or when the game has been
// won. A phase that waits for input is advanced past instead of blocking when an
// effect has ended it early.
func (g *Game) runPhases() {
	for g.State.Winner < 0 && g.State.Phase.valid() {
		if g.State.Phase.waitsForInput() && !g.State.PhaseEnded {
			return
		}
		g.runPhase()
		// A phase can end the game mid-run (forging the third key), so stop here
		// rather than logging entry into a phase that will never actually run.
		if g.State.Winner >= 0 {
			return
		}
		if g.State.Phase == PhaseEndOfTurn {
			return
		}
		g.enterPhase(g.State.Phase + 1)
	}
}

// enterPhase makes p the current phase, clearing the early-end flag that only
// ever applies to the phase that set it.
func (g *Game) enterPhase(p Phase) {
	g.State.Phase = p
	g.State.PhaseEnded = false
	g.record(PhaseBegan{
		Player: g.State.ActivePlayer,
		Phase:  p,
	})
}

// EndPhase ends the current phase early, so the phase loop moves on without
// running the rest of it (Omega ends the play phase the moment it resolves).
func (g *Game) EndPhase() { g.State.PhaseEnded = true }

// runPhase carries out the current phase. An open phase whose body is the
// frontend's own play loop has nothing left to do here; it is only reached once
// an effect has ended it early, and its table entry has no run.
func (g *Game) runPhase() {
	if run := phases[g.State.Phase].run; run != nil {
		run(g, g.State.ActivePlayer)
	}
}

// startOfTurnPhase resolves the active player's "at the start of your turn"
// abilities. It runs before the forge phase, so an ability that changes what a
// key costs still has time to. Its expiry sweep is row-only: the fields it clears
// are per-card row state an upgrade never carries
// (docs/adr/0048-upgrade-in-play-not-an-ability-source.md).
func (g *Game) startOfTurnPhase(player int) {
	// An effect granted "until the start of your next turn" (Hideaway Hole's
	// elusive, the Mutation cycle's Assault and trait, Reckless Rizzo's keyword
	// loss) lasts through the rest of the granting turn and the opponent's whole
	// turn, then lifts here at the start of the controller's next turn — before any
	// start-of-turn ability resolves, so those abilities see the grant already gone.
	for _, id := range g.creaturesAndArtifacts(player) {
		g.State.Cards[id].KeywordsUntilNextTurn = 0
		g.State.Cards[id].LostKeywordsUntilNextTurn = 0
		g.State.Cards[id].AssaultUntilNextTurn = 0
		g.State.Cards[id].TraitUntilNextTurn = traitUnset
	}
	// The active player's cards fire every "at the start of your turn" ability, and
	// the opponent's fire only their EachPlayer-scoped ones — an artifact that
	// watches every turn, not only its owner's (Gambling Den, General Order 24) —
	// each resolving as the active player so "that player" is whoever's turn is
	// starting. The whole set is one window the active player orders (ADR 0013), and
	// it runs before the forge phase so an ability that changes what a key costs
	// still resolves in time.
	start := g.window()
	start.addTurnScoped(player, TriggerStartOfTurn, true)
	g.resolveWindow(g.orderTriggered(player, start.pending))
	g.settleDestroyed(player)
}

// forgePhase forges a key if the player can afford one, unless an effect (Miasma)
// made them skip the phase. The forge itself goes through the Resolver port's
// ForgeKeyAtExtraCost with no surcharge, the same call a forge card makes, so the
// turn's own forge is audited at the port rather than reaching past it.
func (g *Game) forgePhase(player int) {
	if g.State.SkipForge[player].Value || g.skipsForge(player) ||
		g.forgeBarredWhileAhead(player) {
		g.record(ForgeSkipped{Player: player})
		return
	}
	g.ForgeKeyAtExtraCost(player, 0)
}

// drawPhase refills the active player's hand at the end of their turn. It goes
// through the Resolver port's RefillHand, the same call a card that refills a
// hand makes (Punctuated Equilibrium), so the turn's own draw is audited at the
// port rather than reaching past it.
func (g *Game) drawPhase(player int) { g.RefillHand(player) }

// readyPhase readies the active player's cards: it unexhausts every card they
// control and records which ones had been exhausted. It readies and nothing
// else. The state that expires with the turn is lifted in the end-of-turn
// phase's cleanup tail, two phases later, so an end-of-turn ability still sees
// this turn's effects (ADR 0047). Its board scan is row-only: exhaustion is
// per-card row state an upgrade never carries
// (docs/adr/0048-upgrade-in-play-not-an-ability-source.md).
func (g *Game) readyPhase(player int) {
	var readied []LocalID
	for _, id := range g.creaturesAndArtifacts(player) {
		if g.State.Cards[id].Exhausted {
			readied = append(readied, id)
		}
		// Readying goes through the Resolver port's SetExhausted, the same call a
		// card that readies a creature makes, so the turn's own ready is audited at
		// the port rather than reaching past it.
		g.SetExhausted(id, false)
	}
	if len(readied) > 0 {
		g.record(CardsReadied{
			Player: player,
			Cards:  readied,
		})
	}
	// End of the "ready cards" step: every card has readied (Greater Oxtet purges
	// a card from hand to grow).
	for _, id := range g.creaturesAndArtifacts(player) {
		g.triggerAbilities(id, TriggerEndOfReadyStep, 0, false)
	}
	g.settleDestroyed(player)
}

// endOfTurnPhase resolves the active player's "at the end of your turn" abilities
// and the effects scheduled into this turn's end-of-turn window (Ragnarok's board
// wipe), then cleans up everything that expires with the turn. It is the last
// phase, after ready and draw, so those abilities see the board and hand the turn
// actually ends with, and the whole window — in-play abilities and scheduled
// effects alike — is gathered up front and ordered as one (ADR 0013).
//
// The phase runs in two parts. The ability window comes first and the active
// player orders it. The cleanup tail comes second, after the window and its
// settle, so a "remainder of the turn" effect is still in force while the
// abilities that must observe it resolve (ADR 0047). Only the window is
// orderable: the cleanup is bookkeeping with no card-visible interleaving, so it
// runs in whatever fixed order the engine finds convenient and no player is asked
// about it.
func (g *Game) endOfTurnPhase(player int) {
	// The active player's cards fire every "at the end of your turn" ability, and
	// the opponent's fire only their EachPlayer-scoped ones — an artifact that
	// watches every turn, not only its owner's (Pincerator) — each resolving as the
	// active player so "that player" is whoever's turn is ending. The whole set,
	// together with the effects scheduled into this end-of-turn window (Ragnarok's
	// board wipe), is one window the active player orders (ADR 0013).
	w := g.window()
	w.addTurnScoped(player, TriggerEndOfTurn, true)
	w.pending = append(w.pending, g.scheduledEndOfTurn(player)...)
	g.resolveWindow(g.orderTriggered(player, w.pending))
	g.settleDestroyed(player)
	g.expireTurnScoped(player)
	g.clearScheduled()
	// Duration-scoped continuous effects (damage immunity, lost keywords, blanked
	// text, stat overrides) end here, after the end-of-turn window and immediately
	// before the turn hands over, so an end-of-turn ability still sees this turn's
	// effects and a "during your opponent's next turn" effect survives into it.
	g.clearExpiredContinuous()
	// Expiring the turn's effects can leave a creature lethal: a power buff that
	// lasted only the turn, an un-blanked constant power-reducer (Shadow of Dis
	// wearing off King of the Crag), or a lifted stat override can drop a creature
	// to zero power, so settle again before the turn hands over.
	g.settleDestroyed(player)
	// The turn is handed over on a shared scoreboard: the player who just played,
	// then the one about to.
	for _, p := range [2]int{player, 1 - player} {
		g.record(PlayerStanding{
			Player:    p,
			Aember:    g.Aember(p),
			KeyColors: g.KeyColors(p),
		})
	}
	g.assertInvariants()
}

// expireTurnScoped lifts everything that lasts only for the turn now ending: a
// temporary animation (Animator), keywords and stat bonuses granted for the turn
// (Scout, Protectrix), the turn's restriction bars and permits, the key-cost
// bumps, and the lasting-effect registry. It also refreshes creature armor for
// the turn to come and rolls this turn's history into "last turn".
//
// It runs in the end-of-turn phase's cleanup tail, after the end-of-turn
// abilities have resolved, so those abilities still see this turn's effects in
// force (ADR 0047): an end-of-turn ability reaches a creature Animator animated
// this turn as a creature, proven by TestAnimatorRevertsAfterEndOfTurnAbilities.
// The order of the steps below is not player-visible and not orderable.
//
// Unlike readying, forging, and drawing, the cleanup writes state directly rather
// than through the Resolver port, and deliberately stays off it: refreshing armor,
// lifting the turn's bars, and rolling TurnHistory are not capabilities any card
// has, and resolver.go describes the port as the full explicit catalogue of what
// an effect may do — adding them would make that catalogue lie to a card author.
// The cleanup is covered by the witness tests in game_phase_test.go instead, each
// asserting that one step appends the log entry it should.
//
// Its board scans are row-only: exhaustion, armor, and the turn-scoped bonuses it
// clears are per-card row state an upgrade never carries
// (docs/adr/0048-upgrade-in-play-not-an-ability-source.md).
func (g *Game) expireTurnScoped(player int) {
	// A card animated only for this turn (Animator) reverts to an artifact now,
	// keeping its power counters.
	g.revertTemporaryCreatures()
	// "Cannot be dealt damage" lasts only the turn, so clear it on every creature,
	// including any enemy one an effect protected (Protectrix). A keyword gained for
	// the turn (Scout) expires the same way, on whichever creature holds it.
	for _, id := range append(g.creaturesAndArtifacts(player), g.creaturesAndArtifacts(1-player)...) {
		g.State.Cards[id].GrantedKeywords = 0
		g.State.Cards[id].LostKeywords = 0
		g.State.Cards[id].ConsideredFlank = false
		g.State.Cards[id].TempPowerBonus = 0
		g.State.Cards[id].TempArmorBonus = 0
		g.State.Cards[id].TempAssaultBonus = 0
		g.State.Cards[id].TextBoxTurnSourcePlus = 0
	}
	// Armor absorbed or stripped this turn is restored for the turn to come, and a
	// card that belonged to another house for the remainder of the turn goes back to
	// its printed one.
	for _, id := range g.creaturesAndArtifacts(player) {
		core := &g.State.Cards[id]
		core.TempHouse = HouseNone
		if g.TypeOf(id) == Creature {
			core.ArmorRemaining = int16(g.armor(id))
			core.ArmorStripped = 0
		}
	}
	g.State.CannotFight[player] = Bar[bool]{}
	g.State.CannotUse[player] = Bar[bool]{}
	g.State.CannotReap[player] = Bar[bool]{}
	g.State.CannotReapHouse[player] = Bar[House]{}
	g.State.CreaturesCannot[player] = Bar[CreatureBar]{}
	g.State.CannotPlayTypeThis[player] = Bar[CardType]{}
	// Roll this turn's history into "last turn" so the next player can ask what their
	// opponent just did.
	h := &g.State.TurnHistory
	h[player][KeysForgedLastTurn] = h[player][KeysForgedThisTurn]
	h[player][KeysForgedThisTurn] = 0
	h[player][CreaturesPlayedLastTurn] = int8(g.creaturesPlayedThisTurn(player))
	h[0][EnemyCreaturesFightKilled] = 0
	h[1][EnemyCreaturesFightKilled] = 0
	h[0][EnemyCreaturesDestroyed] = 0
	h[1][EnemyCreaturesDestroyed] = 0
	h[0][FriendlyCreaturesDestroyed] = 0
	h[1][FriendlyCreaturesDestroyed] = 0
	h[player][CreaturesReapedThisTurn] = 0
	h[player][CreaturesFoughtThisTurn] = 0
	// The player who acts next reads how much was stolen from them during this
	// just-ended turn as "their opponent's previous turn" (Information Exchange).
	next := 1 - player
	h[next][AemberStolenFromLastTurn] = h[next][AemberStolenFromThisTurn]
	h[next][AemberStolenFromThisTurn] = 0
	h[player][AemberStolenFromThisTurn] = 0
	g.State.MayFightHouse[player] = HouseNone
	g.State.MayFightAny[player] = false
	g.State.MayUseHouse[player] = HouseNone
	g.State.MayPlayHouse[player] = HouseNone
	g.State.MayUseArtifactsAnyHouse[player] = false
	g.State.MayUseTrait[player] = traitUnset
	g.clearOffHousePermits(player)
	g.State.KeyCostBump[player] = Bar[int]{}
	g.State.KeyCostPerHouse[player] = Bar[perHouseKeySurcharge]{}
	g.clearLasting(player)
}
