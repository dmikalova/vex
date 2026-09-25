package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Squawker
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Choose one:
//	- Ready a Mars creature
//	- Stun a non-Mars creature.
var Squawker = set.New(
	"Squawker",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, "178"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.ChooseOne{
			Options: []card.Effect{
				card.Ready{Target: card.Target.Creature.House(card.Houses.Named(card.House.Self))},
				card.Stun{Target: card.Target.Creature.House(card.Houses.Except(card.House.Self))},
			},
		}),
)
