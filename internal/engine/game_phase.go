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
	for _, id := range g.allInPlay(player) {
		g.triggerAbilities(id, TriggerStartOfTurn, 0, false)
	}
	g.settleDestroyed(player)
}

// forgePhase forges a key if the player can afford one, unless an effect (Miasma)
// made them skip the phase.
func (g *Game) forgePhase(player int) {
	if g.State.SkipForge[player].Value {
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
		if g.cat.def(id).Type == Creature {
			core.ArmorRemaining = int16(g.armor(id))
		}
	}
	if len(readied) > 0 {
		g.record(CardsReadied{Player: player, Cards: readied})
	}
	// "Cannot be dealt damage" lasts only the turn, so clear it on every creature,
	// including any enemy one an effect protected (Protectrix).
	for _, id := range append(g.allInPlay(player), g.allInPlay(1-player)...) {
		g.State.Cards[id].DamageImmune = false
	}
	g.State.CannotFight[player] = Bar[bool]{}
	g.State.CannotUse[player] = Bar[bool]{}
	g.State.CannotPlayTypeThis[player] = Bar[CardType]{}
	// Roll this turn's history into "last turn" so the next player can ask what their
	// opponent just did.
	h := &g.State.TurnHistory
	h[player][KeysForgedLastTurn] = h[player][KeysForgedThisTurn]
	h[player][KeysForgedThisTurn] = 0
	h[player][CreaturesPlayedLastTurn] = int8(g.creaturesPlayedThisTurn(player))
	h[0][EnemyCreaturesFightKilled] = 0
	h[1][EnemyCreaturesFightKilled] = 0
	g.State.MayFightHouse[player] = HouseNone
	g.State.MayFightAny[player] = false
	g.State.MayUseHouse[player] = HouseNone
	g.State.KeyCostBump[player] = Bar[int]{}
	g.State.KeywordsLost = 0
	g.clearLasting(player)
	// A power buff that lasted only the turn has just expired.
	g.settleDestroyed(player)
}

// endOfTurnPhase resolves the active player's "at the end of your turn" abilities.
// It is the last phase, after ready and draw, so those abilities see the board and
// hand the turn actually ends with (ADR 0013).
func (g *Game) endOfTurnPhase(player int) {
	for _, id := range g.allInPlay(player) {
		g.triggerAbilities(id, TriggerEndOfTurn, 0, false)
	}
	g.settleDestroyed(player)
	// The turn is handed over on a shared scoreboard: the player who just played,
	// then the one about to.
	for _, p := range [2]int{player, 1 - player} {
		g.record(PlayerStanding{
			Player: p,
			Aember: g.State.Aember[p],
			Keys:   g.State.Keys[p],
		})
	}
	g.assertInvariants()
}
