//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// GrimlocusDux
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Rare
//	Power:  11
//	Armor:  2
//	Traits: Dinosaur • Soldier
//
//	Taunt.
//	Play: Exalt Grimlocus Dux twice.
var GrimlocusDux = card.New(
	"Grimlocus Dux",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, 221),
	card.WithPower(11),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Soldier),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
