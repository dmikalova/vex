package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Professor Sutterkin
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Human • Scientist
//
//	Reap: For each friendly Logos creature, draw a card.
var ProfessorSutterkin = card.New(
	"Professor Sutterkin",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 118),
	card.WithPower(2),
	card.WithTraits(card.Traits.Human, card.Traits.Scientist),
	card.WithAbility(
		card.Trigger.Reap, card.Draw{
			Amount: 1,
			Per: card.InPlay{
				Player: card.Controller,
				Type:   card.Type.Creature,
				House:  card.House.Self,
			},
		}),
)
