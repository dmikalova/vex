//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// SeismoEntangler
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Item
//
//	Action: Choose a house. During your opponent's next turn, creatures of the chosen house cannot be used to reap.
var SeismoEntangler = card.New(
	"Seismo-entangler",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 137),
	card.WithTraits(card.Traits.Item),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
