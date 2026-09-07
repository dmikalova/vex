//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// PrimusUnguis
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Armor:  1
//	Traits: Dinosaur • Soldier
//
//	Each friendly creature gets +2 power for each A on Primus Unguis.
//	Reap: Exalt Primus Unguis.
var PrimusUnguis = card.New(
	"Primus Unguis",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, 226),
	card.WithPower(5),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Soldier),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
