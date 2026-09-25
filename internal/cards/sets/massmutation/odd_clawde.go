package massmutation

import "github.com/dmikalova/vex/internal/card"

// Odd Clawde
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Special
//	Power:  5
//	Traits: Mutant • Scientist
//
//	Action: If your opponent has an odd amount of Æmber, steal 1 Æmber.
var OddClawde = set.New(
	"Odd Clawde",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Special,
	card.Provenance(card.MM, "121"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Mutant, card.Traits.Scientist),
	card.WithAbility(
		card.Trigger.Action, card.Conditional{
			Cond: card.PoolAember{
				Player: card.Opponent,
				Is:     card.Odd,
			},
			Then: card.StealAember{Amount: 1},
		}),
)
