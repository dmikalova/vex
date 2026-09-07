//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// NogiSmartfist
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Traits: Giant • Scientist
//
//	Fight: Draw 2 cards. Discard 2 random cards from your hand.
var NogiSmartfist = card.New(
	"Nogi Smartfist",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, 44),
	card.WithPower(5),
	card.WithTraits(card.Traits.Giant, card.Traits.Scientist),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
