//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// ForgemasterOg
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Giant
//
//	After a player forges a key, they lose all of their remaining A.
var ForgemasterOg = card.New(
	"Forgemaster Og",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 38),
	card.WithPower(4),
	card.WithTraits(card.Traits.Giant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
