//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// EDAIEdie4x4
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Ai • Scientist
//
//	Play: Archive a card.
//	Your opponent's keys cost +1A for each card in your archives.
var EDAIEdie4x4 = card.New(
	"EDAI \"Edie\" 4x4",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, 132),
	card.WithPower(3),
	card.WithTraits(card.Traits.Ai, card.Traits.Scientist),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
