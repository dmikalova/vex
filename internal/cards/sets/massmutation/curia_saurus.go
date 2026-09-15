package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Curia Saurus
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Location
//
//	Each creature with Æmber on it gains, "Destroyed: Move 1 Æmber from this creature to the most powerful enemy creature."
var CuriaSaurus = set.New(
	"Curia Saurus",
	card.House.Saurian,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "202"),
	card.WithTraits(card.Traits.Location),
	card.WithConstant(card.ConstantAbility{
		Target: card.Target.EachCreature.WithAember(),
		Granted: []card.Ability{{
			Trigger: card.Trigger.Destroyed,
			Effect: card.MoveAember{
				Amount: 1,
				From:   card.Target.This,
				Onto:   card.Target.EachEnemyCreature.Refine(card.MostPowerful),
			},
		}},
	}),
)
