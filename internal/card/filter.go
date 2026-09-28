package card

import "github.com/dmikalova/vex/internal/engine"

// The axes a target's filters narrow on, each re-exported as a grouped namespace
// the way card.House and card.Type are: card.Power.AtMost(3) is plainly a power
// bound, card.Position.OnFlank plainly a place in a battleline. Each axis holds
// one value at a time, so a card names one power question, one damage question
// and one place — never two that could contradict each other. Treat these
// package-level vars as read-only.

// Power groups the power filters, e.g. card.Power.AtMost(3). The bounds are
// constructors because each reads a number; the parities are values because
// neither does.
var Power = powerBounds{
	Odd:  engine.PowerBound{Kind: engine.BoundOdd},
	Even: engine.PowerBound{Kind: engine.BoundEven},
}

type powerBounds struct {
	// Odd admits a creature whose power is odd, printed "with odd power" (Onyx
	// Knight).
	Odd engine.PowerBound
	// Even admits a creature whose power is even, printed "with even power" (Opal
	// Knight).
	Even engine.PowerBound
}

// AtMost admits a creature whose power is n or lower, printed "with power n or
// lower".
func (powerBounds) AtMost(n int) engine.PowerBound {
	return engine.PowerBound{Kind: engine.BoundAtMost, Amount: n}
}

// AtLeast admits a creature whose power is n or higher, printed "with power n or
// higher".
func (powerBounds) AtLeast(n int) engine.PowerBound {
	return engine.PowerBound{Kind: engine.BoundAtLeast, Amount: n}
}

// Exactly admits a creature whose power is exactly n, printed "with power n".
func (powerBounds) Exactly(n int) engine.PowerBound {
	return engine.PowerBound{Kind: engine.BoundExactly, Amount: n}
}

// Damage groups the damage filters, e.g. card.Damage.Some.
var Damage = damagePresences{
	Some: engine.DamageSome,
	None: engine.DamageNone,
}

type damagePresences struct {
	// Some admits a creature that has damage on it ("a damaged creature").
	Some engine.DamagePresence
	// None admits a creature that has no damage on it ("an undamaged creature").
	None engine.DamagePresence
}

// Aember groups the Æmber filters, e.g. card.Aember.None.
var Aember = aemberPresences{
	Some: engine.AemberSome,
	None: engine.AemberNone,
}

type aemberPresences struct {
	// Some admits a card with Æmber on it ("a creature with Æmber on it").
	Some engine.AemberPresence
	// None admits a card with no Æmber on it (Draining Touch destroys a creature
	// with no Æmber on it).
	None engine.AemberPresence
}

// Position groups the battleline places, e.g. card.Position.OnFlank. A card names
// one of them: the union of two places is written as a target per place.
var Position = positions{
	OnFlank:    engine.PositionOnFlank,
	NotOnFlank: engine.PositionNotOnFlank,
	Center:     engine.PositionCenter,
	Right:      engine.PositionRightOfSource,
	Left:       engine.PositionLeftOfSource,
}

type positions struct {
	// OnFlank admits a creature on a flank of its battleline. An artifact a target
	// also reaches passes it untouched (Snudge).
	OnFlank engine.Position
	// NotOnFlank admits a creature that is neither its battleline's leftmost nor
	// its rightmost.
	NotOnFlank engine.Position
	// Center admits the creature in the center of its controller's battleline; an
	// even-sized line has no center (Beware the Ides).
	Center engine.Position
	// Right admits a creature standing to the source card's right (Panpaca, Anga).
	Right engine.Position
	// Left admits a creature standing to the source card's left (Panpaca, Jaga).
	Left engine.Position
}

// Except groups the cards a filter can leave out — what "other" means on the
// card that names it, e.g. card.Except.Source.
var Except = exclusions{
	Source: engine.ExcludeSource,
	It:     engine.ExcludeIt,
	Focus:  engine.ExcludeFocus,
}

type exclusions struct {
	// Source leaves out the card the ability is printed on ("each other friendly
	// card").
	Source engine.Exclusion
	// It leaves out the card in context, and nothing when no effect has put one
	// there — "another" for a selection made after a first pick.
	It engine.Exclusion
	// Focus leaves out the card in context when an effect has put one there, and
	// the source card otherwise, so "another creature" never lands on the card it
	// is defined against (Guardian Demon).
	Focus engine.Exclusion
}
