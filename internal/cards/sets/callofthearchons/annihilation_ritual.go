package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Annihilation Ritual
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Power
//
//	Each Creature gains, "Destroyed: Purge this Creature."
var AnnihilationRitual = set.New(
	"Annihilation Ritual",
	card.House.Dis,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "72"),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Power),
	card.WithConstant(card.ConstantAbility{
		Target: card.Target.EachCreature,
		Granted: []card.Ability{{
			Trigger: card.Trigger.Destroyed,
			Effect:  card.PurgeCreature{Target: card.Target.This},
		}},
	}),
)
