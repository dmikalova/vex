package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Hunter or Hunted?
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Remove a ward from a creature. Ward a creature.
var HunterOrHunted = set.New(
	"Hunter or Hunted?",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, "269"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.RemoveWard{Target: card.Target.Creature},
			card.Ward{Target: card.Target.Creature},
		}}),
)
