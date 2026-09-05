//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// EntropicSwirl
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Choose a creature. For each trait that creature has, deal it 2D and gain 1A.
var EntropicSwirl = card.New(
	"Entropic Swirl",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 143),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
