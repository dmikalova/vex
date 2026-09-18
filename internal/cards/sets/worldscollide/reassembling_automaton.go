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
//	If this creature would be destroyed and there is another friendly creature in play, instead fully heal it, exhaust it, and move it to either flank of its controller's battleline.
var ReassemblingAutomaton = set.New(
	"Reassembling Automaton",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "158"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Robot, card.Traits.Experiment),
	card.WithStatic(card.StaticModifier{
		Replaces: card.Replace{
			When: card.Event.Destroyed,
			Cond: card.CardsInPlay{Player: card.Controller, Type: card.Type.Creature, Other: true},
			With: card.Sequence{Effects: []card.Effect{
				card.Heal{Fully: true, Target: card.Target.Triggering},
				card.Exhaust{Target: card.Target.Triggering},
				card.MoveToFlank{Target: card.Target.Triggering},
			}},
		},
	}),
)
