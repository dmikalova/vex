package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Rustgnawer
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Beast • Insect
//
//	Fight: Destroy an artifact. For each Æmber bonus on it, gain 1 Æmber.
var Rustgnawer = card.New(
	"Rustgnawer",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, "330"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Beast, card.Traits.Insect),
	card.WithAbility(
		card.Trigger.Fight, card.Sentences{Effects: []card.Effect{
			card.Destroy{Target: card.Target.Artifact},
			card.GainAember{
				Player: card.Controller,
				Amount: 1,
				Per:    card.AemberBonusOf{Target: card.Target.Triggering},
			},
		}}),
)
