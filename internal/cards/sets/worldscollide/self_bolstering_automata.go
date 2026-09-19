package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Self-Bolstering Automata
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Robot
//
//	If this creature would be destroyed and there is another friendly creature in play, instead fully heal it, exhaust it, move it to either flank of its controller's battleline, and give it two +1 power counters.
var SelfBolsteringAutomata = set.New(
	"Self-Bolstering Automata",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "176"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Robot),
	card.WithStatic(card.StaticModifier{
		Replaces: card.Replace{
			When: card.Event.Destroyed,
			Cond: card.CardsInPlay{
				Player: card.Controller,
				Type:   card.Type.Creature,
				Other:  true,
			},
			With: card.Sequence{Effects: []card.Effect{
				card.Heal{
					Fully:  true,
					Target: card.Target.Triggering,
				},
				card.Exhaust{Target: card.Target.Triggering},
				card.MoveToFlank{Target: card.Target.Triggering},
				card.AddPowerCounter{
					Target: card.Target.Triggering,
					Amount: 2,
				},
			}},
		},
	}),
)
