package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Nogi Smartfist
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
	card.WithAbility(
		card.Trigger.Fight, card.Sentences{Effects: []card.Effect{
			card.Draw{Amount: 2},
			card.DiscardRandomFromHand{Player: card.Controller, Amount: 2},
		}}),
)
