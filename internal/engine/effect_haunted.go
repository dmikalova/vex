package engine

// Haunted is met while the controller has 10 or more cards in their discard pile.
// A haunted player has no inherent restriction; cards reference the state, e.g.
// The Grim Reaper enters play ready while its controller is haunted.
type Haunted struct{}

// CondText renders the condition as the card prints it.
func (Haunted) CondText() string { return "if you are haunted" }

// Met reports whether the controller's discard pile holds at least 10 cards.
func (Haunted) Met(ctx *EffectContext) bool {
	return len(ctx.Resolver.Discard(ctx.Controller)) >= 10
}
