//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// ZYXResearcher
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Human • Scientist
//
//	Play: Archive the top card of your deck or the top card of your discard pile.
var ZYXResearcher = card.New(
	"Z.Y.X. Researcher",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 123),
	card.WithPower(2),
	card.WithTraits(card.Traits.Human, card.Traits.Scientist),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
