package massmutation

import "github.com/dmikalova/vex/internal/card"

// Savage Clash
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Destroy each enemy creature except the most powerful enemy creature and each friendly creature except the least powerful friendly creature.
var SavageClash = set.New(
	"Savage Clash",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "376"),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.Destroy{
				Target: card.Target.EachEnemyCreature.Refine(card.Except(card.MostPowerful)),
			},
			card.Destroy{
				Target: card.Target.EachFriendlyCreature.Refine(card.Except(card.LeastPowerful)),
			},
		}}),
)
