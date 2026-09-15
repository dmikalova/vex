//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Eunoia
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Traits: Beast • Cat
//
//	After an enemy creature is destroyed fighting Eunoia, gain 1A and heal 2 damage from Eunoia.
var Eunoia = set.New(
	"Eunoia",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "400"),
	card.WithPower(6),
	card.WithTraits(card.Traits.Beast, card.Traits.Cat),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
