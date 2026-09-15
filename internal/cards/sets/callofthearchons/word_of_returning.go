package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Word of Returning
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Deal 1 damage to each enemy creature for each Æmber on it, and move all Æmber from each enemy creature to your pool.
var WordOfReturning = set.New(
	"Word of Returning",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "339"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.DealDamage{
				Amount:    1,
				Target:    card.Target.EachEnemyCreature,
				PerTarget: card.AemberOnIt,
			},
			card.MoveAember{
				All:  true,
				From: card.Target.EachEnemyCreature,
				To:   card.Controller,
			},
		}}),
)
