//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Paraguardian
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  6
//	Armor:  1
//	Traits: Dinosaur • Soldier
//
//	Reap: You may exalt Paraguardian. If you do, ward each of its neighbors.
var Paraguardian = card.New(
	"Paraguardian",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 206),
	card.WithPower(6),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Soldier),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
