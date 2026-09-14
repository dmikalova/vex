package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Igon the Terrible
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Connected
//	Power:  8
//	Traits: Giant
//
//	Play: If Igon the Green has not been purged, destroy Igon the Terrible.
//	Fight: Steal 1 Æmber.
var IgonTheTerrible = set.New(
	"Igon the Terrible",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Connected,
	card.Provenance(card.WC, "53"),
	card.InCluster(igonCluster),
	card.WithPower(8),
	card.WithTraits(card.Traits.Giant),
	card.WithAbility(
		card.Trigger.Play, card.Conditional{
			Cond: card.Not{Cond: card.NamedCardPurged{Name: IgonTheGreenName}},
			Then: card.Destroy{Target: card.Target.This},
		}),
	card.WithAbility(card.Trigger.Fight, card.StealAember{Amount: 1}),
)
