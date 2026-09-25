package massmutation

import "github.com/dmikalova/vex/internal/card"

// Dark Centurion
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Traits: Mutant • Soldier
//
//	Action: Move 1 Æmber from a creature to the common supply -> ward the chosen creature.
//	Enhance Capture Capture.
var DarkCenturion = set.New(
	"Dark Centurion",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "203"),
	card.WithEnhance(card.Bonus.Capture, card.Bonus.Capture),
	card.WithPower(5),
	card.WithTraits(card.Traits.Mutant, card.Traits.Soldier),
	card.WithAbility(
		card.Trigger.Action, card.Then{
			First: card.MoveAemberToSupply{
				Amount: 1,
				Target: card.Target.Creature,
				Bind:   true,
			},
			Result: card.Ward{Target: card.Target.TheChosenCreature},
		}),
)
