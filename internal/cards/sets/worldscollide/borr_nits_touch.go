package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Borr Nit's Touch
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Reveal the top 5 cards of a player's deck. Purge a card revealed this way. Shuffle that deck.
var BorrNitsTouch = card.New(
	"Borr Nit's Touch",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "087"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play,
		card.RevealTopOfDeck{
			Amount:          5,
			ChooseWhoseDeck: true,
			Then: []card.TopAct{
				card.ChooseAndMove{Count: 1, Dest: card.Into.Purge},
				card.ShuffleDeck{},
			},
		},
	),
)
