package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Questor Jarta
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Dinosaur • Politician
//
//	Elusive.
//	Reap: You may exalt Questor Jarta. Gain 1 Æmber.
var QuestorJarta = set.New(
	"Questor Jarta",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "191"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Politician),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Reap, card.May{Do: card.Sequence{Effects: []card.Effect{
			card.Exalt{
				Target: card.Target.This,
				Amount: 1,
			},
			card.GainAember{
				Player: card.Controller,
				Amount: 1,
			},
		}}}),
)
