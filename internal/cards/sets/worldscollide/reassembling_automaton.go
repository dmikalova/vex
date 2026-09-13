package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Reassembling Automaton
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Robot • Experiment
//
//	Destroyed: If you have any other Creatures in play, instead of destroying Reassembling Automaton, fully heal it, exhaust it, and move it to either flank of its controller's battleline.
var ReassemblingAutomaton = card.New(
	"Reassembling Automaton",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "158"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Robot, card.Traits.Experiment),
	card.WithAbility(card.Trigger.Destroyed, card.Conditional{
		Cond: card.HasOtherFriendlyCreatures{},
		Then: card.SaveFromDestruction{
			Do: card.Sequence{Effects: []card.Effect{
				card.Heal{Fully: true, Target: card.Target.Triggering},
				card.Exhaust{Target: card.Target.Triggering},
				card.MoveToFlank{Target: card.Target.Triggering},
			}},
		},
	}),
)
