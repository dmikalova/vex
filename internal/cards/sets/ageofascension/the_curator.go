//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// TheCurator
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Human • Scientist
//
//	Friendly artifacts enter play ready.
var TheCurator = card.New(
	"The Curator",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 157),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human, card.Traits.Scientist),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
