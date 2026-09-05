//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// GuardDisguise
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: Sacrifice Guard Disguise. If your opponent has 3A or fewer, steal 3A.
var GuardDisguise = card.New(
	"Guard Disguise",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 302),
	card.WithTraits(card.Traits.Item),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
