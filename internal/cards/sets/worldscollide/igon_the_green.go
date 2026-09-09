package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Igon the Green
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Giant
//
//	Destroyed: Purge Igon the Green, and put an Igon the Terrible from your discard pile into your hand.
var IgonTheGreen = card.New(
	"Igon the Green",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "39"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Giant),
	card.Connects(card.PullExact(IgonTheTerrible, 1)),
	card.WithAbility(
		card.Trigger.Destroyed, card.Sequence{Effects: []card.Effect{
			card.PurgeCreature{Target: card.Target.This},
			card.PutFromDiscard{
				Name:        "Igon the Terrible",
				Destination: card.To.Hand,
			},
		}}),
)
