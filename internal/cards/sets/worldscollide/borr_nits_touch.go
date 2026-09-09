package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// BorrNitsTouch
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
	card.Provenance(card.WC, "087"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play,
		card.RevealPurgeShuffleDeck{Amount: 5},
	),
)
