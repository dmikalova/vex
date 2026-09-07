//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// CodeMonkey
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Ai • Beast
//
//	Deploy. (This creature can enter play anywhere in your battleline.)
//	Play: Archive each neighboring creature. If those creatures share a house, gain 2A.
var CodeMonkey = card.New(
	"Code Monkey",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "147"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Ai, card.Traits.Beast),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
