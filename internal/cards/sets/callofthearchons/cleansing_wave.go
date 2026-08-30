package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Cleansing Wave
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Heal 1 damage from each creature, and for each creature healed this way, gain 1 Æmber.
var CleansingWave = card.New(
	"Cleansing Wave",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 215),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.Heal{
				Amount: 1,
				Target: card.Target.EachCreature,
			},
			card.GainAember{
				Player: card.Controller,
				Amount: 1,
				Per:    card.CreaturesHealed{},
			},
		}}),
)
