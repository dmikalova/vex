package engine

// SourceIsFighting is met while the source creature is one of the two combatants
// of the fight resolving right now. It gates a combat-scoped self-grant: Nizak,
// The Forgotten gains invulnerable "while fighting", so its ConstantAbility is
// suspended until it enters a fight and lifts again the moment the fight ends.
type SourceIsFighting struct{}

// CondText renders the condition. It carries the "if " prefix every condition
// leads with; the constant-ability renderer strips it and re-frames the whole
// grant as "While fighting, <self> gains ...".
func (SourceIsFighting) CondText() string { return "if fighting" }

// Met reports whether the source creature is currently in a fight.
func (SourceIsFighting) Met(ctx *EffectContext) bool {
	return ctx.Resolver.CurrentlyFighting(ctx.Source)
}
