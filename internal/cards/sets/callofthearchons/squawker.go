package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Squawker
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Choose one:
//	- Ready a Mars creature
//	- Stun a non-Mars creature.
var Squawker = card.New(
	"Squawker",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 178),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.ChooseOne{
			Options: []card.Effect{
				card.Ready{Target: card.Target.Creature.OfHouse(card.House.Self)},
				card.Stun{Target: card.Target.Creature.ExceptHouse(card.House.Self)},
			},
		}),
)
