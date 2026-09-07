package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Eldest Bear
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Traits: Beast • Leader • Witch
//
//	Assault 3.
//	While Eldest Bear is in the center of your battleline, it gains, "Before Fight: Gain 2 Æmber."
var EldestBear = card.New(
	"Eldest Bear",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, 388),
	card.WithPower(5),
	card.WithTraits(card.Traits.Beast, card.Traits.Leader, card.Traits.Witch),
	card.WithAssault(3),
	card.WithConstant(card.ConstantAbility{
		Target:        card.Target.This,
		WhileInCenter: true,
		Granted: []card.Ability{{
			Trigger: card.Trigger.BeforeFight,
			Effect: card.GainAember{
				Player: card.Controller,
				Amount: 2,
			},
		}},
	}),
)
