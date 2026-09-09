package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Manchego
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Human • Thief
//
//	Play: If you have 5 or fewer cards in your deck, steal 2 Æmber.
//	Fight/Reap: You may shuffle Manchego into its owner's deck.
var Manchego = card.New(
	"Manchego",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "275"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human, card.Traits.Thief),
	card.WithAbility(
		card.Trigger.Play, card.Conditional{
			Cond: card.CardsInDeckAtMost{Amount: 5},
			Then: card.StealAember{Amount: 2},
		}),
	card.WithFightOrReap(card.May{Do: card.PutFromPlay{
		Target:      card.Target.This,
		Destination: card.To.DeckShuffled,
	}}),
)
