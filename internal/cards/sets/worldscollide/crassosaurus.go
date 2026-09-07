//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Crassosaurus
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Armor:  2
//	Traits: Dinosaur • Politician
//
//	Elusive.
//	Play: Capture 10A from any combination of players. Then, if Crassosaurus has fewer than 10A on it, purge Crassosaurus.
var Crassosaurus = card.New(
	"Crassosaurus",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "217"),
	card.WithPower(4),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Politician),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
