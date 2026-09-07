//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Spartasaur
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Armor:  1
//	Traits: Dinosaur • Soldier
//
//	After a friendly creature is destroyed, destroy each non-Dinosaur creature.
//	Fight: Gain 2A.
var Spartasaur = card.New(
	"Spartasaur",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "231"),
	card.WithPower(6),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Soldier),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
