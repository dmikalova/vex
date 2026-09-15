//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// PurifierOfSouls
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Armor:  2
//	Traits: Human • Priest
//
//	Destroyed effects cannot trigger.
var PurifierOfSouls = set.New(
	"Purifier of Souls",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "174"),
	card.WithPower(5),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Human, card.Traits.Priest),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
