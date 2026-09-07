//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// BorrNitsTouch
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Reveal the top 5 cards of a player's deck. Purge a card revealed this way. Shuffle the other revealed cards into that deck.
var BorrNitsTouch = card.New(
	"Borr Nit's Touch",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "87"),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
