package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Charge!
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: For the remainder of the turn, each time you play a creature, deal 2 damage to an enemy creature.
var Charge = set.New(
	"Charge!",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "214"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.ForRemainderOfTurn{
			On: card.Event.CreaturePlayed,
			Do: card.DealDamage{
				Amount: 2,
				Target: card.Target.EnemyCreature,
			},
		}),
)
