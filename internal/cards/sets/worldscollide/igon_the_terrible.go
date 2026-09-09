package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Igon the Terrible
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Special
//	Power:  8
//	Traits: Giant
//
//	Play: If Igon the Green has not been purged, destroy Igon the Terrible.
//	Fight: Steal 1 Æmber.
var IgonTheTerrible = card.New(
	"Igon the Terrible",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Special,
	card.Provenance(card.WC, "53"),
	card.WithPower(8),
	card.WithTraits(card.Traits.Giant),
	card.WithAbility(
		card.Trigger.Play, card.Conditional{
			Cond: card.NamedCardPurged{Name: "Igon the Green", Not: true},
			Then: card.Destroy{Target: card.Target.This},
		}),
	card.WithAbility(card.Trigger.Fight, card.StealAember{Amount: 1}),
)
