//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// DracoPraeco
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Dinosaur • Politician
//
//	Reap: You may exalt Draco Praeco. If you do, choose a house. Enrage each creature of that house.
var DracoPraeco = card.New(
	"Draco Praeco",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 201),
	card.WithPower(4),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Politician),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
