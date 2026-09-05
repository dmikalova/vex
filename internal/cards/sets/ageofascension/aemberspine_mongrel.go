//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// AemberspineMongrel
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Beast
//
//	Hazardous 3. (Before this creature is attacked, deal 3D to the attacking enemy.)
//	After your opponent uses a creature to reap, gain 1A.
var AemberspineMongrel = card.New(
	"Aemberspine Mongrel",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 335),
	card.WithPower(3),
	card.WithTraits(card.Traits.Beast),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
