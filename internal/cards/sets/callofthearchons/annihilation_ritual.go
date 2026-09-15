package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Annihilation Ritual
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Power
//
//	Each creature gains, "Destroyed: Purge this creature."
var AnnihilationRitual = set.New(
	"Annihilation Ritual",
	card.House.Dis,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "72"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Power),
	card.WithConstant(card.ConstantAbility{
		Target: card.Target.EachCreature,
		Granted: []card.Ability{{
			Trigger: card.Trigger.Destroyed,
			Effect:  card.PurgeCreature{Target: card.Target.This},
		}},
	}),
)
