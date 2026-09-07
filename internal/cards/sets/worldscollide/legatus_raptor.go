//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// LegatusRaptor
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Armor:  1
//	Traits: Dinosaur • Soldier
//
//	Fight: You may exalt Legatus Raptor. If you do, ready and use another friendly creature.
var LegatusRaptor = card.New(
	"Legatus Raptor",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, 187),
	card.WithPower(4),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Soldier),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
