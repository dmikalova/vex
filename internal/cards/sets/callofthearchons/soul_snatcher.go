package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Soul Snatcher
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Vehicle
//
//	Each creature gains, "Destroyed: Gain 1 Æmber."
var SoulSnatcher = card.New(
	"Soul Snatcher",
	card.House.Dis,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 80),
	card.WithTraits(card.Traits.Vehicle),
	card.WithConstant(card.ConstantAbility{
		Target: card.Target.EachCreature,
		Granted: []card.Ability{{
			Trigger: card.Trigger.Destroyed,
			Effect: card.GainAember{
				Player: card.Controller,
				Amount: 1,
			},
		}},
	}),
)
