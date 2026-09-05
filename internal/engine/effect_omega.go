package engine

// This file holds the Omega keyword: a card with Omega ends the current step of
// the turn the moment it resolves. The play gates in game_play.go call
// endStepIfOmega after a successful play from hand.

// endStepIfOmega ends the current step of the turn when the active player has just
// played an Omega card. Omega says no more cards may be played, used, or discarded
// for the rest of the step, except through pending abilities still resolving — the
// card's own play effect (and anything it triggers) has already fully resolved by
// the time this runs, so those pending plays still happen and play then continues
// to the next step (Unlocked Gateway, Little Niff). It is a no-op unless the card
// carries Omega and the player is still mid-play-phase, so an effect-driven play
// never trips it.
func (g *Game) endStepIfOmega(player int, def *CardDefinition) {
	if def.hasKeyword(Omega) &&
		g.State.Winner < 0 &&
		g.State.Phase == PhasePlay &&
		g.State.ActivePlayer == player {
		g.EndPlayPhase(player)
	}
}
