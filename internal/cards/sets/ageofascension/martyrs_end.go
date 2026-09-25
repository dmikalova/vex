package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Martyr's End
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Destroy any number of friendly creatures. For each creature destroyed this way, gain 1 Æmber.
var MartyrsEnd = set.New(
	"Martyr's End",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "255"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.DestroyChosen{Target: card.Target.EachFriendlyCreature},
			card.GainAember{
				Player: card.Controller,
				Amount: 1,
				Per:    card.CreaturesDestroyed{},
			},
		}}),
)
