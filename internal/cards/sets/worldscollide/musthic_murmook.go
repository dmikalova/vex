//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// MusthicMurmook
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Beast
//
//	Each player's keys cost +1A.
//	Play: Deal 4D to a creature.
var MusthicMurmook = card.New(
	"Musthic Murmook",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, 361),
	card.WithPower(4),
	card.WithTraits(card.Traits.Beast),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
