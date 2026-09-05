//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// OnyxKnight
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Demon • Knight
//
//	Play: Destroy each creature with odd power.
var OnyxKnight = card.New(
	"Onyx Knight",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 95),
	card.WithPower(4),
	card.WithTraits(card.Traits.Demon, card.Traits.Knight),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
