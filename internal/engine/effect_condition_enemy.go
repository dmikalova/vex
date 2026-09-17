package engine

// ItIsEnemy is met when the card in context (ctx.It — a triggering card) is
// controlled by the ability's opponent. It filters a board-wide trigger, which
// fires for both players' cards, down to the opponent's creatures — the mirror of
// ItIsFriendly — so "after an enemy creature reaps" is the board-wide reap trigger
// narrowed to the creatures you do not control.
type ItIsEnemy struct{}

// CondText renders the condition, e.g. "if it is an enemy creature".
func (ItIsEnemy) CondText() string { return "if it is an enemy creature" }

// Met reports whether a card is in context and its controller is not the
// ability's controller.
func (ItIsEnemy) Met(ctx *EffectContext) bool {
	return ctx.HasIt && ctx.Resolver.Controller(ctx.It) != ctx.Controller
}
