package engine

import "fmt"

// PlayRequirement is an Æmber requirement a card puts on its own play: the pool
// its controller must have before the card can leave their hand, and whether
// playing it spends that Æmber. Two cards word the same gate differently —
// Kelifi Dragon only checks the pool, Truebaru charges it — so one value covers
// both, and the zero value requires nothing.
//
// It is a flat comparable value that lives on the card definition, never in game
// state, so it costs nothing to carry on every card.
type PlayRequirement struct {
	// Aember is how much the controller must have in their pool.
	Aember int
	// Spend charges the requirement rather than only checking it.
	Spend bool
}

// AemberThreshold requires a pool of at least n to play the card, spending
// nothing (Kelifi Dragon).
func AemberThreshold(n int) PlayRequirement { return PlayRequirement{Aember: n} }

// AemberCost requires — and spends — n Æmber from the pool to play the card
// (Truebaru).
func AemberCost(n int) PlayRequirement { return PlayRequirement{Aember: n, Spend: true} }

// required reports whether the card carries a requirement at all.
func (r PlayRequirement) required() bool { return r.Aember > 0 }

// met reports whether a pool satisfies the requirement.
func (r PlayRequirement) met(pool int) bool { return pool >= r.Aember }

// text renders the requirement as the card prints it; SelfName stands in for the
// card's name, which cardRules substitutes.
func (r PlayRequirement) text() string {
	if !r.required() {
		return ""
	}
	if r.Spend {
		return fmt.Sprintf("You must lose %d Æmber in order to play %s.", r.Aember, SelfName)
	}
	return fmt.Sprintf(
		"%s cannot be played unless you have %d Æmber or more.",
		SelfName, r.Aember,
	)
}
