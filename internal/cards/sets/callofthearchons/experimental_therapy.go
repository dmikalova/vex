package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Experimental Therapy
//
//	House:  Logos
//	Type:   Upgrade
//	Rarity: Uncommon
//
//	This Creature gains versatile.
//	Play: Stun and exhaust this Creature.
var ExperimentalTherapy = set.New(
	"Experimental Therapy",
	card.House.Logos,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "157"),
	card.WithStatic(card.StaticModifier{Keywords: card.Keywords(card.Keyword.Versatile)}),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.Stun{Target: card.Target.This},
			card.Exhaust{Target: card.Target.This},
		}}),
)
