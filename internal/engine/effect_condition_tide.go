package engine

// This file holds the tide conditions: the "if the tide is low/high" checks a card
// reads. They resolve against the ability's controller, so "the tide is low" means
// low for the player whose card is checking. The tide has no raiser yet (see
// tide.go), so both read false in every current game — enough to ship a card whose
// only tide interaction is a check (Valoocanth is barred while the tide is low).

// TideIsLow is met when the tide is low for the ability's controller.
type TideIsLow struct{}

// CondText renders the condition, e.g. "While the tide is low, ...".
func (TideIsLow) CondText() string { return "if the tide is low" }

// Met reports whether the tide is low for the ability's controller.
func (TideIsLow) Met(ctx *EffectContext) bool {
	return ctx.Resolver.TideIsLow(ctx.Controller)
}

// TideIsHigh is met when the tide is high for the ability's controller.
type TideIsHigh struct{}

// CondText renders the condition, e.g. "if the tide is high".
func (TideIsHigh) CondText() string { return "if the tide is high" }

// Met reports whether the tide is high for the ability's controller.
func (TideIsHigh) Met(ctx *EffectContext) bool {
	return ctx.Resolver.TideIsHigh(ctx.Controller)
}
