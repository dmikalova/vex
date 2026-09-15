package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Forge Compiler
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Item
//
//	After your opponent forges a key, destroy Forge Compiler, and ward each friendly creature.
var ForgeCompiler = set.New(
	"Forge Compiler",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "088"),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.AfterOpponentForgesKey, card.Sequence{
			Effects: []card.Effect{
				card.Destroy{Target: card.Target.This},
				card.Ward{Target: card.Target.EachFriendlyCreature},
			},
		}),
)
