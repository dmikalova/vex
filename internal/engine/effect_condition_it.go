package engine

// ItIsOfHouse is met when the card in context (ctx.It — a revealed, discarded, or
// triggering card) belongs to a referenced house. It replaces the one-off
// "revealed card of the chosen house" (Chaos Portal) and "discarded card of the
// active house" (Evasion Sigil) with a single filter on the contextual card.
type ItIsOfHouse struct {
	House HouseChoice
}

// CondText renders the condition, e.g. "if it is of the chosen house".
func (e ItIsOfHouse) CondText() string {
	if e.House == TheActiveHouse {
		return "if it is of the active house"
	}
	return "if it is of the chosen house"
}

// Met reports whether a card is in context and belongs to the referenced house.
func (e ItIsOfHouse) Met(ctx *EffectContext) bool {
	if !ctx.HasIt {
		return false
	}
	house := ctx.Resolver.House(ctx.It)
	if e.House == TheActiveHouse {
		return house == ctx.Resolver.ActiveHouse()
	}
	return house == ctx.ChosenHouse
}

// ItIs is met when the card in context (ctx.It — a just-played, revealed, or
// discarded card) matches a concrete House and/or Type filter, e.g. "if it is a
// Mars creature" (Brain Stem Antenna reacting to a played card) or "if it is an
// artifact" (Carlo Phantom). Either filter may be left unset to match any. It is
// the concrete-value counterpart to ItIsOfHouse, which names the house by
// reference (the chosen or active house).
type ItIs struct {
	// House and Type are the filters the card in context must match; either one
	// left unset matches any.
	House House
	Type  CardType
	// Not inverts the match, so the condition is met when the card in context does
	// NOT fit the filters — Neutron Shark repeats until it discards a Logos card.
	Not bool
	// Other excludes the source card itself, so "it" must be a different card and
	// the noun reads "another" — Hunting Witch gains only when you play another
	// creature, never on its own entrance (Harmonia, which says "a creature", omits
	// it and gains from its own play).
	Other bool
	// Subject names the card outright when "it" has drifted too far from the trigger
	// that set it. Unset says "it".
	Subject Subject
}

// shapeNoun renders the house/type shape the contextual card must match, prefixed
// "another" when Other bars the source card itself.
func (e ItIs) shapeNoun() string {
	if e.Other {
		return "another " + houseTypeNoun(e.House, e.Type)
	}
	return houseTypeNoun(e.House, e.Type)
}

// CondText renders the condition, e.g. "if it is a Mars creature", "if it is an
// artifact", or, inverted and named, "if the discarded card is not a Logos card".
func (e ItIs) CondText() string {
	if e.Not {
		return "if " + e.Subject.noun() + " is not " + indefinite(e.shapeNoun())
	}
	return "if " + e.Subject.noun() + " is " + indefinite(e.shapeNoun())
}

// Met reports whether a card is in context and matches the house and type
// filters, inverting the match under Not. Other additionally bars the source card
// itself, so a card never counts its own play.
func (e ItIs) Met(ctx *EffectContext) bool {
	if !ctx.HasIt {
		return false
	}
	if e.Other && ctx.It == ctx.Source {
		return false
	}
	return e.matches(ctx) != e.Not
}

// matches reports whether the card in context fits the house and type filters.
func (e ItIs) matches(ctx *EffectContext) bool {
	if e.House != HouseNone && ctx.Resolver.House(ctx.It) != e.House {
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
