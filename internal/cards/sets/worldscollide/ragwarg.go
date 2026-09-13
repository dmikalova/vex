package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Ragwarg
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Item
//
//	After a Creature reaps, if it is the first time a Creature has reaped this turn, deal 2 damage to it.
var Ragwarg = card.New(
	"Ragwarg",
	card.House.Brobnar,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "27"),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.AfterCreatureReaps, card.Conditional{
			Cond: card.FirstReapOfTurn{},
			Then: card.DealDamage{
				Target: card.Target.Triggering,
				Amount: 2,
			},
		}),
)
