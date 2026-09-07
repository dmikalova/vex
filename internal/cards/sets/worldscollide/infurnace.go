//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Infurnace
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Demon
//
//	Play: Purge up to 2 cards from a discard pile. Your opponent loses A equal to the total Aember bonus of the purged cards.
var Infurnace = card.New(
	"Infurnace",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, 78),
	card.WithPower(4),
	card.WithTraits(card.Traits.Demon),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
