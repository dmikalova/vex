package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Irradiated Aember
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: If your opponent has 6 Æmber or more, deal 3 damage to each enemy creature.
var IrradiatedAember = set.New(
	"Irradiated Aember",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "165"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Conditional{
			Cond: card.PoolAember{
				Player: card.Opponent,
				Is:     card.AtLeast,
				Amount: 6,
			},
			Then: card.DealDamage{
				Amount: 3,
				Target: card.Target.EachEnemyCreature,
			},
		}),
)
