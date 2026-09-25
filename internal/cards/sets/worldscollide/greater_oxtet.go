package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Greater Oxtet
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Demon
//
//	Taunt.
//	At the end of your "ready cards" step, purge a card from your hand -> give Greater Oxtet two +1 power counters.
var GreaterOxtet = set.New(
	"Greater Oxtet",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "105"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Demon),
	card.WithKeywords(card.Keyword.Taunt),
	card.WithAbility(
		card.Trigger.EndOfReadyStep, card.Then{
			First: card.PurgeCard{
				Zones:     []card.Zone{card.Hand},
				Player:    card.Controller,
				Selection: card.Chosen{},
			},
			Result: card.AddPowerCounter{
				Target: card.Target.This,
				Amount: 2,
			},
		}),
)
