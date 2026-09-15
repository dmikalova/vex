package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Hayden Oswin
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Human
//
//	Reap: For each upgrade on Hayden Oswin, gain 1 Æmber.
var HaydenOswin = set.New(
	"Hayden Oswin",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "322"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Human),
	card.WithAbility(
		card.Trigger.Reap, card.GainAember{
			Player: card.Controller,
			Amount: 1,
			Per:    card.UpgradesOn{Target: card.Target.This},
		}),
)
