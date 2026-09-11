package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Sanitation Engineer
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Cyborg • Scientist
//
//	Hazardous 1.
//	Reap: Discard a card from your hand.
var SanitationEngineer = card.New(
	"Sanitation Engineer",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "138"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Cyborg, card.Traits.Scientist),
	card.WithHazardous(1),
	card.WithAbility(
		card.Trigger.Reap, card.DiscardCard{
			Player:    card.Controller,
			Zone:      card.Hand,
			Selection: card.Chosen{},
			Amount:    1,
		}),
)
