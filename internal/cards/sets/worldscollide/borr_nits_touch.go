package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Borr Nit's Touch
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Reveal the top 5 cards of a player's deck. Purge a card revealed this way. Shuffle that deck.
var BorrNitsTouch = set.New(
	"Borr Nit's Touch",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "087"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play,
		card.RevealTopOfDeck{
			Amount:          5,
			ChooseWhoseDeck: true,
			Then: []card.TopAct{
				card.ChooseAndMove{
					Cards: 1,
					Dest:  card.Into.Purge,
				},
				card.Shuffle{},
			},
		},
	),
)
