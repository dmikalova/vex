package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Munchling
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Mutant
//
//	Skirmish.
//	Fight: You may discard a Logos card from your hand or archives -> gain 1 Æmber.
var Munchling = set.New(
	"Munchling",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "076"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Mutant),
	card.WithKeywords(card.Keyword.Skirmish),
	card.WithAbility(card.Trigger.Fight, card.May{
		Do: card.Then{
			First: card.DiscardCard{
				Player:    card.Controller,
				Zones:     []card.Zone{card.Hand, card.Archives},
				Selection: card.Chosen{House: card.Houses.Named(card.House.Self)},
			},
			Result: card.GainAember{
				Player: card.Controller,
				Amount: 1,
			},
		},
	}),
)
