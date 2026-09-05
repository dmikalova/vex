//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// AubadeTheGrim
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Armor:  1
//	Traits: Spirit • Knight
//
//	Play: Capture 3A.
//	Reap: Discard 1A from Aubade the Grim.
var AubadeTheGrim = card.New(
	"Aubade the Grim",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 213),
	card.WithPower(4),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Spirit, card.Traits.Knight),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
