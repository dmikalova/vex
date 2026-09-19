package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Borr Nit
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Demon
//
//	Reap: Reveal the top 5 cards of a player's deck. Purge a card revealed this way. Shuffle that deck.
var BorrNit = set.New(
	"Borr Nit",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "086"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Demon),
	card.WithAbility(
		card.Trigger.Reap,
		card.RevealTopOfDeck{
			Amount:          5,
			ChooseWhoseDeck: true,
			Then: []card.TopAct{
				card.ChooseAndMove{Cards: 1, Dest: card.Into.Purge},
				card.Shuffle{},
			},
		},
	),
)
