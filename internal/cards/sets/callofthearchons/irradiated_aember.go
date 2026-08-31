package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Irradiated Aember
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: If your opponent has 6 Æmber or more, deal 3 damage to each enemy creature.
var IrradiatedAember = card.New(
	"Irradiated Aember",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 165),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Conditional{
			Cond: card.OpponentAember{Is: card.AtLeast, Amount: 6},
			Then: card.DealDamage{
				Amount: 3,
				Target: card.Target.EachEnemyCreature,
			},
		}),
)
