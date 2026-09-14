package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Thorium Plasmate
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Move an enemy Creature anywhere in its controller's battleline -> for each neighbor that shares a house with the chosen Creature, deal 2 damage to the chosen Creature.
var ThoriumPlasmate = set.New(
	"Thorium Plasmate",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "140"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Then{
			First: card.MoveWithinBattleline{Target: card.Target.EnemyCreature},
			Result: card.DealDamage{
				Amount: 2,
				Per:    card.NeighborsSharingHouse{},
				Target: card.Target.TheChosenCreature,
			},
		}),
)
