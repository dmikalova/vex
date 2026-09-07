package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Helmsman Spears
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Human
//
//	Fight/Reap: Discard any number of cards from your hand -> for each card discarded this way, draw a card.
var HelmsmanSpears = card.New(
	"Helmsman Spears",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 311),
	card.WithPower(2),
	card.WithTraits(card.Traits.Human),
	card.WithFightOrReap(card.Then{
		First:  card.DiscardFromHand{AnyNumber: true},
		Result: card.ForEachDiscarded{Do: card.Draw{Amount: 1}},
	}),
)
