//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// SenatorBracchus
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Dinosaur • Politician
//
//	You may spend A on friendly creatures as if it were in your pool.
//	Fight/Reap: Exalt Senator Bracchus.
var SenatorBracchus = card.New(
	"Senator Bracchus",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "229"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Politician),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
