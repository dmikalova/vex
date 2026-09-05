package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Carpet Phloxem
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: If there are no friendly creatures in play, deal 4 damage to each creature.
var CarpetPhloxem = card.New(
	"Carpet Phloxem",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, 161),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Conditional{
			Cond: card.InPlay{
				Player: card.Controller,
				Type:   card.Type.Creature,
				None:   true,
			},
			Then: card.DealDamage{
				Amount: 4,
				Target: card.Target.EachCreature,
			},
		}),
)
