package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Carpet Phloxem
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: If there are no friendly creatures in play, deal 4 damage to each creature.
var CarpetPhloxem = set.New(
	"Carpet Phloxem",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, "161"),
	card.WithBonus(card.Bonus.Aember),
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
