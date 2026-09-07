//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// BorrNit
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Demon
//
//	Reap: Reveal the top 5 cards of a player's deck. Purge a card revealed this way. Shuffle the other revealed cards into that deck.
var BorrNit = card.New(
	"Borr Nit",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 86),
	card.WithPower(3),
	card.WithTraits(card.Traits.Demon),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
