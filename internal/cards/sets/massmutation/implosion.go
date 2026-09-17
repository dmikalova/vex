package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Imp-losion
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Destroy a friendly creature and an enemy creature.
var Implosion = set.New(
	"Imp-losion",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "008"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.BatchDestroy{
			Gather: card.ChosenFromEach{
				card.Target.FriendlyCreature,
				card.Target.EnemyCreature,
			},
		}),
)
