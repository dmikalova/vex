package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Phloxem Spike
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: If there are no friendly creatures in play, destroy each creature that is not on a flank.
var PhloxemSpike = card.New(
	"Phloxem Spike",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 186),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Conditional{
			Cond: card.InPlay{
				Player: card.Controller,
				Type:   card.Type.Creature,
				None:   true,
			},
			Then: card.Destroy{Target: card.Target.EachCreature.NotOnFlank()},
		}),
)
