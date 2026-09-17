package engine

// ItIsAmong is met when the subject creature (ctx.It) is among the creatures a
// Target selects — the general "is it one of these" test. Subject names the card
// "it" refers to in the printed text (the fought creature, that card), defaulting
// to a bare "it". Baldric the Bold reads it in a Before Fight ability to gain
// Æmber when the creature it fights is the most powerful enemy creature; a tie
// still qualifies, and asking the membership never prompts.
type ItIsAmong struct {
	Target  Target
	Subject Subject
}

// CondText renders the condition with the "if" prefix every condition leads with.
func (c ItIsAmong) CondText() string {
	return "if " + c.Subject.noun() + " is " + c.Target.Text()
}

// Met reports whether the subject creature is among the Target's tie-inclusive
// candidates, without making the discretionary tie-break the Target's own
// selection would.
func (c ItIsAmong) Met(ctx *EffectContext) bool {
	if !ctx.HasIt {
		return false
	}
	return c.Target.couldSelect(ctx, ctx.It)
}
