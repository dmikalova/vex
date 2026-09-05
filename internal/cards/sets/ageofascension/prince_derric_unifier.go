//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// PrinceDerricUnifier
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Armor:  1
//	Traits: Human • Knight
//
//	Play: Gain 3A if you control creatures from 3 different houses.
var PrinceDerricUnifier = card.New(
	"Prince Derric, Unifier",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 240),
	card.WithPower(4),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Human, card.Traits.Knight),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
