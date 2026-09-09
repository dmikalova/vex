package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// ReassemblingAutomaton
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Robot • Experiment
//
//	Destroyed: If you have any other creatures in play, instead of destroying Reassembling Automaton, fully heal it, exhaust it, and move it to a flank.
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
