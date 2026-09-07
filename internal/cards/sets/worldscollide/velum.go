//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Velum
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Special
//	Power:  2
//	Traits: Human • Scientist
//
//	Reap: Archive a card. If you control Hyde, archive 2 cards instead.
//	Destroyed: Archive Hyde from your discard pile. If you do, archive Velum.
var Velum = card.New(
	"Velum",
	card.House.Logos,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from FIXED to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.WC, "181"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Human, card.Traits.Scientist),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
