package engine

// ItIs is met when the card in context (ctx.It — a just-played, revealed, or
// discarded card) matches a House and/or Type filter, e.g. "if it is a Mars
// creature" (Brain Stem Antenna reacting to a played card), "if it is an artifact"
// (Carlo Phantom), or "if it is of the chosen house" (Chaos Portal). The House
// matcher is the one house-filter vocabulary (ADR 0038): a named house, a
// non-<house> ("if it is a non-Star Alliance card" — Book of leQ), or a referenced
// house (the chosen or active house). Either filter may be left unset to match any.
type ItIs struct {
	// House and Type are the filters the card in context must match; either one
	// left unset (the zero HouseMatcher / TypeUnset) matches any.
	House HouseMatcher
	Type  CardType
	// Other excludes the source card itself, so "it" must be a different card and
	// the noun reads "another" — Hunting Witch gains only when you play another
	// creature, never on its own entrance (Harmonia, which says "a creature", omits
	// it and gains from its own play).
	Other bool
	// Subject names the card outright when "it" has drifted too far from the trigger
	// that set it. Unset says "it".
	Subject Subject
}

// shapeNoun renders the house/type shape the contextual card must match for a
// prefix-kind matcher, e.g. "Mars creature" or "non-Logos card", prefixed
// "another" when Other bars the source card itself.
func (e ItIs) shapeNoun() string {
	noun := e.House.qualifyNoun(typeNoun(e.Type))
	if e.Other {
		return "another " + noun
	}
	return noun
}

// predicate renders what the card in context must be. A referenced house reads as
// a trailing phrase ("of the chosen house"); every other matcher reads as an
// article-and-noun ("a Mars creature", "a non-Logos card").
func (e ItIs) predicate() string {
	switch e.House.Kind {
	case MatchChosenHouse:
		return "of the chosen house"
	case MatchActiveHouse:
		return "of the active house"
	default:
		return indefinite(e.shapeNoun())
	}
}

// CondText renders the condition, e.g. "if it is a Mars creature" or "if it is of
// the chosen house".
func (e ItIs) CondText() string {
	return "if " + e.Subject.noun() + " is " + e.predicate()
}

// negatedText renders the inverted clause a Not wrapper prints, e.g. "if the
// discarded card is not a Logos card" (Neutron Shark).
func (e ItIs) negatedText() string {
	return "if " + e.Subject.noun() + " is not " + e.predicate()
}

// Met reports whether a card is in context and matches the house and type
// filters. Other additionally bars the source card itself, so a card never counts
// its own play.
func (e ItIs) Met(ctx *EffectContext) bool {
	if !ctx.HasIt {
		return false
	}
	if e.Other && ctx.It == ctx.Source {
		return false
	}
	return e.matches(ctx)
}

// matches reports whether the card in context fits the house and type filters.
func (e ItIs) matches(ctx *EffectContext) bool {
	if !e.House.matches(ctx, ctx.It) {
		return false
	}
	if e.Type != TypeUnset && ctx.Resolver.TypeOf(ctx.It) != e.Type {
		return false
	}
	return true
}

// ItIsOfTrait is met when the creature in context (ctx.It) has the named trait.
type ItIsOfTrait struct{ Trait Trait }

// CondText renders the condition, e.g. "if it is a Dinosaur creature".
func (c ItIsOfTrait) CondText() string {
	return "if it is a " + c.Trait.String() + " creature"
}

// Met reports whether a creature is in context and has the trait.
func (c ItIsOfTrait) Met(ctx *EffectContext) bool {
	return ctx.HasIt && ctx.Resolver.HasTrait(ctx.It, c.Trait)
}

// ItHasAember is met when the creature in context (ctx.It) has any Æmber on it.
type ItHasAember struct{}

// CondText renders the condition.
func (ItHasAember) CondText() string { return "if it has \u00c6mber on it" }

// Met reports whether a creature is in context with Æmber on it.
func (ItHasAember) Met(ctx *EffectContext) bool {
	return ctx.HasIt && ctx.Resolver.AmberOn(ctx.It) > 0
}

// ItHasBonusIcon is met when the card in context (ctx.It — a just-played card) has
// at least one printed bonus icon, the gate on Adaptoid's "after you play a card
// with a bonus icon" reaction.
type ItHasBonusIcon struct{}

// CondText renders the condition.
func (ItHasBonusIcon) CondText() string { return "if it has a bonus icon" }

// Met reports whether a card is in context and prints at least one bonus icon.
func (ItHasBonusIcon) Met(ctx *EffectContext) bool {
	return ctx.HasIt && ctx.Resolver.HasBonusIcons(ctx.It)
}

// ItIsOffIdentity is met when the card in context (ctx.It) belongs to none of the
// controller's identity houses — the three houses of their deck. Sneklifter uses
// it to reassign a seized enemy artifact to Shadows only when it is off your
// identity; KeyForge phrases the check negatively ("if it does not belong to one
// of your three houses").
type ItIsOffIdentity struct{}

// CondText renders the condition.
func (ItIsOffIdentity) CondText() string {
	return "if it does not belong to a house on your identity"
}

// Met reports whether a card is in context and belongs to none of the controller's
// identity houses.
func (ItIsOffIdentity) Met(ctx *EffectContext) bool {
	return ctx.HasIt && !ctx.Resolver.PlayerHasHouse(ctx.Controller, ctx.Resolver.House(ctx.It))
}

// ItIsStunned is met when the creature in context (ctx.It) is already stunned.
// 1-2 Punch uses it to destroy a chosen creature that was already stunned rather
// than stunning it.
type ItIsStunned struct{}

// CondText renders the condition.
func (ItIsStunned) CondText() string {
	return "if that creature was already stunned"
}

// Met reports whether a creature is in context and is stunned.
func (ItIsStunned) Met(ctx *EffectContext) bool {
	return ctx.HasIt && ctx.Resolver.Stunned(ctx.It)
}

// ItIsNotOfNamedHouse is met when a card is in context (ctx.It) and is not of the
// house a player named earlier in this ability, stored in ctx.ChosenHouse by an
// OpponentNamesHouse. Keyforgery uses it: a revealed card not of the named house
// destroys the guard and cancels the forge. Like ItIsOffIdentity it renders a
// negative sentence but is a positive, HasIt-gated condition, so an empty hand
// (no card revealed) leaves it unmet and the forge proceeds. It is its own
// condition, not an ItIs house matcher, because the named house is dynamic and the
// HouseMatcher facade's Named selector is the fixed-house one.
type ItIsNotOfNamedHouse struct {
	// Subject names the card in context outright; unset says "it".
	Subject Subject
}

// CondText renders the condition, e.g. "if that card is not of the named house".
func (c ItIsNotOfNamedHouse) CondText() string {
	return "if " + c.Subject.noun() + " is not of the named house"
}

// Met reports whether a card is in context and is not of the named house.
func (c ItIsNotOfNamedHouse) Met(ctx *EffectContext) bool {
	return ctx.HasIt && ctx.Resolver.House(ctx.It) != ctx.ChosenHouse
}
