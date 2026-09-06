package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Martyr's End
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Destroy any number of friendly creatures. For each creature destroyed this way, gain 1 Æmber.
var MartyrsEnd = card.New(
	"Martyr's End",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 255),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{Effects: []card.Effect{
			card.DestroyChosen{Target: card.Target.EachFriendlyCreature},
			card.GainAember{
				Player: card.Controller,
				Amount: 1,
				Per:    card.CreaturesDestroyed{},
			},
		}}),
)
