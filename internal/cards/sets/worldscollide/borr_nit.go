package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// BorrNit
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
	card.Provenance(card.WC, "086"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Demon),
	card.WithAbility(
		card.Trigger.Reap,
		card.RevealPurgeShuffleDeck{Amount: 5},
	),
)
