package massmutation

import "github.com/dmikalova/vex/internal/card"

// Defense Initiative
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Choose a creature. Ward the chosen creature. You may exalt the chosen creature -> ward each neighbor of the chosen creature.
var DefenseInitiative = set.New(
	"Defense Initiative",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "191"),
	card.WithAbility(
		card.Trigger.Play, card.ChooseCreatureThen{
			Target: card.Target.Creature,
			Then: card.Sequence{Effects: []card.Effect{
				card.Ward{Target: card.Target.TheChosenCreature},
				card.May{Do: card.Then{
					First: card.Exalt{
						Target: card.Target.TheChosenCreature,
						Amount: 1,
					},
					Result: card.Ward{Target: card.Target.TheChosenCreature.NeighborsOf()},
				}},
			}},
		}),
)
