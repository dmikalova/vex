package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Z-Force Agent 14
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Traits: Cyborg
//
//	Fight: For each upgrade on Z-Force Agent 14, gain 1 Æmber.
var ZForceAgent14 = set.New(
	"Z-Force Agent 14",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "353"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Cyborg),
	card.WithAbility(
		card.Trigger.Fight, card.GainAember{
			Player: card.Controller,
			Amount: 1,
			Per:    card.UpgradesOn{Target: card.Target.This},
		}),
)
