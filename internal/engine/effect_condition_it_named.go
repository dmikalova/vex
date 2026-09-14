package engine

// ItIsNamed is met when the card in context (ctx.It — a just-played, revealed, or
// triggering card) carries the given printed name — Chain Gang readies only when
// the card you played is Subtle Chain. It is the by-name counterpart to ItIs,
// which filters the contextual card by house and type.
type ItIsNamed struct {
	Name string
}

// CondText renders the condition, e.g. "if it is Subtle Chain".
func (e ItIsNamed) CondText() string {
	return "if it is " + e.Name
}

// Met reports whether a card is in context and carries the given name.
func (e ItIsNamed) Met(ctx *EffectContext) bool {
	return ctx.HasIt && CardFilter{Name: e.Name}.admits(ctx.Resolver, ctx.It)
}
