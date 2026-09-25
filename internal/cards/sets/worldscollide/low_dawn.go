package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Low Dawn
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: If there are 3 or more Untamed creatures in your discard pile, gain 2 Æmber. Shuffle each Untamed creature from your discard pile into your deck.
var LowDawn = set.New(
	"Low Dawn",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "377"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(card.Trigger.Play, card.Sequence{Effects: []card.Effect{
		card.Conditional{
			Cond: card.CardsInDiscardAtLeast{
				House:  card.Houses.Named(card.House.Self),
				Type:   card.Type.Creature,
				Amount: 3,
			},
			Then: card.GainAember{
				Player: card.Controller,
				Amount: 2,
			},
		},
		card.ShuffleIntoDeck{
			Player: card.Controller, From: []card.Zone{card.Discard}, Selection: card.Each{
				House: card.Houses.Named(card.House.Self),
				Type:  card.Type.Creature,
			}},
	}}),
)
