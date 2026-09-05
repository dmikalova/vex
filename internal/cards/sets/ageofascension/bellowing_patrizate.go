//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// BellowingPatrizate
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  7
//	Traits: Giant
//
//	While Bellowing Patrizate is ready, each creature takes 1D after it enters play.
var BellowingPatrizate = card.New(
	"Bellowing Patrizate",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 34),
	card.WithPower(7),
	card.WithTraits(card.Traits.Giant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
