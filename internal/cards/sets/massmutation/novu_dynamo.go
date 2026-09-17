package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Novu Dynamo
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  8
//	Armor:  2
//	Traits: Robot
//
//	At the start of your turn, discard a Logos card from your hand or archives -> gain 1 Æmber. Otherwise, destroy Novu Dynamo.
var NovuDynamo = set.New(
	"Novu Dynamo",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "093"),
	card.WithPower(8),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Robot),
	card.WithAbility(
		card.Trigger.StartOfTurn, card.Then{
			First: card.DiscardCard{
				Player: card.Controller,
				Zones:  []card.Zone{card.Hand, card.Archives},
				Selection: card.Chosen{
					House:    card.Houses.Named(card.House.Self),
					Optional: true,
				},
			},
			Result: card.GainAember{
				Player: card.Controller,
				Amount: 1,
			},
			Else: card.Destroy{Target: card.Target.This},
		}),
)
