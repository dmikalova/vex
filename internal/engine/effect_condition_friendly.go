package engine

// ItIsFriendly is met when the card in context (ctx.It — a just-played or
// triggering card) is controlled by the ability's controller. It filters the
// global "after a creature is played" trigger, which fires for both players'
// cards, down to the controller's own creatures — Sci. Officer Morpheus
// re-triggers the play effect only of a creature you played.
type ItIsFriendly struct{}

// CondText renders the condition, e.g. "if it is a friendly creature".
func (ItIsFriendly) CondText() string { return "if it is a friendly creature" }

// Met reports whether a card is in context and its controller is the ability's
// controller.
func (ItIsFriendly) Met(ctx *EffectContext) bool {
	return ctx.HasIt && ctx.Resolver.Controller(ctx.It) == ctx.Controller
}
