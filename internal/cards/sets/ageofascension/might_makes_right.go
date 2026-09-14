package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Might Makes Right
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Destroy any number of friendly Creatures. If the total power of Creatures destroyed this way is 25 or more, forge a key at no cost -> purge Might Makes Right.
var MightMakesRight = set.New(
	"Might Makes Right",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "43"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{Effects: []card.Effect{
			card.DestroyChosen{Target: card.Target.EachFriendlyCreature},
			card.Conditional{
				Cond: card.CountIs{
					Count:  card.PowerDestroyedThisWay{},
					Is:     card.AtLeast,
					Amount: 25,
				},
				Then: card.ForgeKey{FreeOfCost: true},
			},
		}}),
)
