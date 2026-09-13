package engine

// This file holds the whole-hand refresh effects: DiscardHand discards a player's
// entire hand, RefillHand redraws it as if the turn had ended. Punctuated
// Equilibrium composes both over EachPlayer.

// playersInOrder lists the players a Player value acts on: the single named player,
// or, for EachPlayer, both players with the active player first — the active player
// orders simultaneous per-player effects.
func playersInOrder(ctx *EffectContext, p Player) []int {
	if p == EachPlayer {
		active := ctx.Resolver.ActivePlayer()
		return []int{active, 1 - active}
	}
	return []int{ctx.PlayerFor(p)}
}

// DiscardHand discards a player's whole hand, one card at a time so a card with a
// discard reaction sees each. Player may be EachPlayer, discarding both hands with
// the active player's first.
type DiscardHand struct {
	Player Player
}

// validate rejects a DiscardHand whose player was left unset.
func (e DiscardHand) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("DiscardHand")
	}
	return nil
}

// Text renders the effect, e.g. "each player discards their hand".
func (e DiscardHand) Text() string {
	switch e.Player {
	case Opponent:
		return "your opponent discards their hand"
	case EachPlayer:
		return "each player discards their hand"
	default:
		return "discard your hand"
	}
}

// Resolve discards each affected player's whole hand, one card at a time.
func (e DiscardHand) Resolve(ctx *EffectContext) {
	for _, p := range playersInOrder(ctx, e.Player) {
		for _, id := range ctx.Resolver.Hand(p) {
			ctx.Resolver.DiscardCardFromHand(p, id)
		}
	}
}

// RefillHand has a player refill their hand as if it were the end of their turn,
// honoring their chains and draw modifiers. Player may be EachPlayer, refilling
// both hands with the active player's first.
type RefillHand struct {
	Player Player
}

// validate rejects a RefillHand whose player was left unset.
func (e RefillHand) validate() error {
	if !e.Player.valid() {
		return errUnsetPlayer("RefillHand")
	}
	return nil
}

// Text renders the effect, e.g. "each player refills their hand as if it were the
// end of their turn".
func (e RefillHand) Text() string {
	switch e.Player {
	case Opponent:
		return "your opponent refills their hand as if it were the end of their turn"
	case EachPlayer:
		return "each player refills their hand as if it were the end of their turn"
	default:
		return "refill your hand as if it were the end of the turn"
	}
}

// Resolve refills each affected player's hand as if their turn had ended.
func (e RefillHand) Resolve(ctx *EffectContext) {
	for _, p := range playersInOrder(ctx, e.Player) {
		ctx.Resolver.RefillHand(p)
	}
}
