package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Hunter or Hunted?
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Choose one:
//	- Ward a creature
//	- Move a ward from a creature to another creature.
var HunterOrHunted = card.New(
	"Hunter or Hunted?",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, "269"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.ChooseOne{
			Options: []card.Effect{
				card.Ward{Target: card.Target.Creature},
				card.MoveWard{
					From: card.Target.Creature,
					Onto: card.Target.OtherCreature,
				},
			},
		}),
)
