package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Even Ivan
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Mutant • Scientist
//
//	Action: If your opponent has an even amount of Æmber, steal 1 Æmber.
var EvenIvan = set.New(
	"Even Ivan",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "073"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Mutant, card.Traits.Scientist),
	card.WithAbility(
		card.Trigger.Action, card.Conditional{
			Cond: card.PoolAember{Player: card.Opponent, Is: card.Even},
			Then: card.StealAember{Amount: 1},
		}),
)
