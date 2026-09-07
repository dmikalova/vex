//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// CincinnatusRex
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Armor:  4
//	Traits: Dinosaur • Soldier
//
//	If there are no enemy creatures, destroy Cincinnatus Rex.
//	Fight: You may exalt Cincinnatus Rex. If you do, ready each other friendly card.
var CincinnatusRex = card.New(
	"Cincinnatus Rex",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, 215),
	card.WithPower(6),
	card.WithArmor(4),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Soldier),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
