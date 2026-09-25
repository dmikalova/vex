package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Thorium Plasmate
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Move an enemy creature anywhere in its controller's battleline -> for each neighbor of that card's house, deal 2 damage to the chosen creature.
var ThoriumPlasmate = set.New(
	"Thorium Plasmate",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "140"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Then{
			First: card.MoveWithinBattleline{Target: card.Target.EnemyCreature},
			Result: card.DealDamage{
				Amount: 2,
				Per:    card.NeighborsMatching{House: card.Houses.Contextual},
				Target: card.Target.TheChosenCreature,
			},
		}),
)
