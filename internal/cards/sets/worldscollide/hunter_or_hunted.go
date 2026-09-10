package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Hunter or Hunted?
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Remove a ward from a creature, and ward a creature.
var HunterOrHunted = card.New(
	"Hunter or Hunted?",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, "269"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.RemoveWard{Target: card.Target.Creature},
			card.Ward{Target: card.Target.Creature},
		}}),
)
