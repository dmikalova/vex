package engine

// EndTurn ends the active player's turn the moment it resolves — the turn simply
// stops, with no ready, draw, or end-of-turn abilities. Book of leQ ends your turn
// this way when the card it reveals is a Star Alliance card. This is unlike the
// Omega keyword, which ends the play phase but still runs the turn out (ready,
// draw, and end-of-turn abilities); EndTurn skips all of that cleanup, so the
// active player's cards stay exhausted and their hand is not refilled.
type EndTurn struct{}

// Text renders the effect.
func (EndTurn) Text() string { return "end your turn" }

// Resolve ends the active player's turn.
func (EndTurn) Resolve(ctx *EffectContext) { ctx.Resolver.EndTurnNow() }

// EndTurnNow ends the active player's turn immediately, jumping straight to the
// end-of-turn phase without running the ready, draw, or end-of-turn steps. The
// turn stops where it stands: cards stay exhausted, no cards are drawn, and no
// end-of-turn abilities fire. This is the opposite of the Omega keyword, which
// ends the play phase but still runs those steps (EndPlayPhase). It is a no-op
// once the game is won or the play phase has already ended, so ending an
// already-ending turn does nothing.
func (g *Game) EndTurnNow() {
	if g.State.Winner < 0 && g.State.Phase == PhasePlay {
		g.enterPhase(PhaseEndOfTurn)
	}
}
