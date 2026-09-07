//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Daughter
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Cyborg • Scientist
//
//	Elusive. (The first time this creature is attacked each turn, no damage is dealt.)
//	During your "draw cards" step, refill your hand to 1 additional card.
var Daughter = card.New(
	"Daughter",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "131"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Cyborg, card.Traits.Scientist),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
