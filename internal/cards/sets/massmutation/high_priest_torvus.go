//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// HighPriestTorvus
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Armor:  1
//	Traits: Dinosaur • Priest
//
//	Reap: You may exalt High Priest Torvus. If you do, after you resolve your next action card this turn, return it to your hand instead of placing it in your discard pile.
var HighPriestTorvus = set.New(
	"High Priest Torvus",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "222"),
	card.WithPower(4),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Priest),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
