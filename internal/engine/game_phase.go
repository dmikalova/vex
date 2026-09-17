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
	g.record(PhaseBegan{Player: g.State.ActivePlayer, Phase: p})
}

// EndPhase ends the current phase early, so the phase loop moves on without
// running the rest of it (Omega ends the play phase the moment it resolves).
func (g *Game) EndPhase() { g.State.PhaseEnded = true }

// runPhase carries out the current phase. An open phase whose body is the
// frontend's own play loop has nothing left to do here; it is only reached once
// an effect has ended it early.
func (g *Game) runPhase() {
	player := g.State.ActivePlayer
	switch g.State.Phase {
	case PhaseStartOfTurn:
		g.startOfTurnPhase(player)
	case PhaseForge:
		g.forgePhase(player)
	case PhaseArchives:
		g.offerArchives(player)
	case PhaseReady:
		g.readyPhase(player)
	case PhaseDraw:
		g.drawStep(player)
	case PhaseEndOfTurn:
		g.endOfTurnPhase(player)
	}
}

// startOfTurnPhase resolves the active player's "at the start of your turn"
// abilities. It runs before the forge phase, so an ability that changes what a
// key costs still has time to.
func (g *Game) startOfTurnPhase(player int) {
	// An effect granted "until the start of your next turn" (Hideaway Hole's
	// elusive, the Mutation cycle's Assault and trait, Reckless Rizzo's keyword
	// loss) lasts through the rest of the granting turn and the opponent's whole
	// turn, then lifts here at the start of the controller's next turn — before any
	// start-of-turn ability resolves, so those abilities see the grant already gone.
	for _, id := range g.allInPlay(player) {
		g.State.Cards[id].KeywordsUntilNextTurn = 0
		g.State.Cards[id].LostKeywordsUntilNextTurn = 0
		g.State.Cards[id].AssaultUntilNextTurn = 0
		g.State.Cards[id].TraitUntilNextTurn = traitUnset
	}
	for _, id := range g.allInPlay(player) {
		g.triggerAbilities(id, TriggerStartOfTurn, 0, false)
	}
	// A start-of-turn artifact watches every turn, not only its owner's (Gambling
	// Den, General Order 24). This window fires for both players' in-play cards, each
	// resolving as the turn's active player so "they"/"that player" is the player
	// whose turn is starting, not the artifact's controller.
	for _, p := range [2]int{player, 1 - player} {
		for _, id := range g.allInPlay(p) {
			g.triggerAbilitiesAs(player, id, TriggerAfterAnyPlayerStartOfTurn, 0, false)
		}
	}
	g.settleDestroyed(player)
}

// forgePhase forges a key if the player can afford one, unless an effect (Miasma)
// made them skip the phase.
func (g *Game) forgePhase(player int) {
	if g.State.SkipForge[player].Value || g.skipsForge(player) ||
		g.forgeBarredWhileAhead(player) {
		g.record(ForgeSkipped{Player: player})
		return
	}
	g.forgeKey(player)
}

// readyPhase readies the active player's cards, refreshes creature armor for the
// new turn, and lifts the this-turn-only state that expires with the turn.
func (g *Game) readyPhase(player int) {
	var readied []LocalID
	for _, id := range g.allInPlay(player) {
		core := &g.State.Cards[id]
		if core.Exhausted {
			readied = append(readied, id)
		}
		core.Exhausted = false
		core.TempHouse = HouseNone
		if g.TypeOf(id) == Creature {
			core.ArmorRemaining = int16(g.armor(id))
			core.ArmorStripped = 0
		}
	}
	if len(readied) > 0 {
		g.record(CardsReadied{Player: player, Cards: readied})
	}
	// A card animated only for this turn (Animator) reverts to an artifact now, at
	// end of turn, keeping its power counters.
	g.revertTemporaryCreatures()
	// "Cannot be dealt damage" lasts only the turn, so clear it on every creature,
	// including any enemy one an effect protected (Protectrix). A keyword gained for
	// the turn (Scout) expires the same way, on whichever creature holds it.
	for _, id := range append(g.allInPlay(player), g.allInPlay(1-player)...) {
		g.State.Cards[id].GrantedKeywords = 0
		g.State.Cards[id].LostKeywords = 0
		g.State.Cards[id].ConsideredFlank = false
		g.State.Cards[id].TempPowerBonus = 0
		g.State.Cards[id].TempArmorBonus = 0
		g.State.Cards[id].TempAssaultBonus = 0
		g.State.Cards[id].TextBoxTurnSourcePlus = 0
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
	// End of the "ready cards" step: every card has readied (Greater Oxtet purges
	// a card from hand to grow).
	for _, id := range g.allInPlay(player) {
		g.triggerAbilities(id, TriggerEndOfReadyStep, 0, false)
	}
	// A power buff that lasted only the turn has just expired.
	g.settleDestroyed(player)
}

// endOfTurnPhase resolves the active player's "at the end of your turn" abilities
// and the effects scheduled into this turn's end-of-turn window (Ragnarok's board
// wipe). It is the last phase, after ready and draw, so those abilities see the
// board and hand the turn actually ends with, and the whole window — in-play
// abilities and scheduled effects alike — is gathered up front and ordered as one
// (ADR 0013).
func (g *Game) endOfTurnPhase(player int) {
	var pending []triggeredAbility
	for _, id := range g.allInPlay(player) {
		abilities := g.triggeredBy(id, TriggerEndOfTurn)
		for i := range abilities {
			abilities[i].actor = int8(player)
		}
		pending = append(pending, abilities...)
	}
	// An end-of-turn artifact watches every turn, not only its owner's (Pincerator).
	// This window fires for both players' in-play cards, each resolving as the turn's
	// active player so "they"/"that player" is the player whose turn is ending, not
	// the artifact's controller.
	for _, p := range [2]int{player, 1 - player} {
		for _, id := range g.allInPlay(p) {
			abilities := g.triggeredBy(id, TriggerAfterAnyPlayerEndOfTurn)
			for i := range abilities {
				abilities[i].actor = int8(player)
			}
			pending = append(pending, abilities...)
		}
	}
	pending = append(pending, g.scheduledEndOfTurn(player)...)
	g.resolveWindow(g.orderTriggered(player, pending))
	g.settleDestroyed(player)
	g.clearScheduled()
	// Duration-scoped continuous effects (damage immunity, lost keywords, blanked
	// text, stat overrides) end here, after the end-of-turn window and immediately
	// before the turn hands over, so an end-of-turn ability still sees this turn's
	// effects and a "during your opponent's next turn" effect survives into it.
	g.clearExpiredContinuous()
	// Expiring a continuous effect can leave a creature lethal: an un-blanked
	// constant power-reducer (Shadow of Dis wearing off King of the Crag) or a
	// lifted stat override can drop a creature to zero power, so settle again
	// before the turn hands over.
	g.settleDestroyed(player)
	// The turn is handed over on a shared scoreboard: the player who just played,
	// then the one about to.
	for _, p := range [2]int{player, 1 - player} {
		g.record(PlayerStanding{
			Player:    p,
			Aember:    g.State.Aember[p],
			KeyColors: g.KeyColors(p),
		})
	}
	g.assertInvariants()
}
