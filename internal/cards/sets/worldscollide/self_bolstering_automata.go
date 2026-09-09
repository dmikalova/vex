package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// SelfBolsteringAutomata
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Robot
//
//	Destroyed: If you have any other creatures in play, instead of destroying Self-Bolstering Automata, fully heal it, exhaust it, and move it to a flank. If you do, give it two +1 power counters.
var SelfBolsteringAutomata = card.New(
	"Self-Bolstering Automata",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "176"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Robot),
	card.WithAbility(card.Trigger.Destroyed, card.Conditional{
		Cond: card.HasOtherFriendlyCreatures{},
		Then: card.Then{
			First: card.SaveFromDestruction{
				Do: card.Sequence{Effects: []card.Effect{
					card.Heal{Fully: true, Target: card.Target.Triggering},
					card.Exhaust{Target: card.Target.Triggering},
					card.MoveToFlank{Target: card.Target.Triggering},
				}},
			},
			Result: card.AddPowerCounter{Target: card.Target.Triggering, Amount: 2},
		},
	}),
)
