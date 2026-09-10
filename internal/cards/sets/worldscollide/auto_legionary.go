package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Auto-Legionary
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Robot • Ally
//
//	Versatile.
//	Action: Give Auto-Legionary five +1 power counters. Move it to a flank of your battleline as a creature.
var AutoLegionary = card.New(
	"Auto-Legionary",
	card.House.Saurian,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "214"),
	card.WithTraits(card.Traits.Robot, card.Traits.Ally),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.Action, card.Sentences{Effects: []card.Effect{
			card.AddPowerCounter{Target: card.Target.This, Amount: 5},
			card.TurnIntoCreature{Target: card.Target.This},
		}}),
)
