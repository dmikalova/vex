//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// MercyMalkinQueen
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Human • Witch
//
//	Skirmish.
//	After a friendly Cat creature enters play, ward it.
//	Fight: Ready a friendly Beast creature.
var MercyMalkinQueen = set.New(
	"Mercy, Malkin Queen",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "403"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human, card.Traits.Witch),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
