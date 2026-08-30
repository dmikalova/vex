package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Nerve Blast
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Steal 1 Æmber -> deal 2 damage to a creature.
var NerveBlast = card.New(
	"Nerve Blast",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 276),
	card.WithAbility(
		card.Trigger.Play, card.Then{
			First: card.StealAember{Amount: 1},
			Result: card.DealDamage{
				Amount: 2,
				Target: card.Target.Creature,
			},
		}),
)
