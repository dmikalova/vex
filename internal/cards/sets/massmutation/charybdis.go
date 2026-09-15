//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Charybdis
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Special
//	Power:  7
//	Traits: Beast
//
//	Each enemy creatures gains, "Before Fight: Lose 1A."
var Charybdis = set.New(
	"Charybdis",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Special,
	card.Provenance(card.MM, "234"),
	card.WithPower(7),
	card.WithTraits(card.Traits.Beast),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
