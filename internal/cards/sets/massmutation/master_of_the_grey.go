//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// MasterOfTheGrey
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Armor:  1
//	Traits: Human • Monk
//
//	Your opponent cannot resolve bonus icons on cards they play.
var MasterOfTheGrey = set.New(
	"Master of the Grey",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "169"),
	card.WithPower(4),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Human, card.Traits.Monk),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
