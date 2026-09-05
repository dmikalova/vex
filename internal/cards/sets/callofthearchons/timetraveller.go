package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Timetraveller
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Human • Scientist
//
//	Play: Draw 2 cards.
//	Action: Shuffle Timetraveller into its owner's deck.
var Timetraveller = card.New(
	"Timetraveller",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 153),
	card.WithPower(2),
	card.WithTraits(card.Traits.Human, card.Traits.Scientist),
	card.Connects(
		card.Pull(HelpFromFutureSelf, 1),
	),
	card.WithAbility(
		card.Trigger.Play, card.Draw{Amount: 2}),
	card.WithAbility(
		card.Trigger.Action, card.PutFromPlay{
			Target:      card.Target.This,
			Destination: card.To.DeckShuffled,
		}),
)
