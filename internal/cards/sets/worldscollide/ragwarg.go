package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Ragwarg
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Uncommon
//	Bonus:  Æmber
//	Traits: Item
//
//	After a creature reaps, if this is the first time a creature has reaped this turn, deal 2 damage to it.
var Ragwarg = set.New(
	"Ragwarg",
	card.House.Brobnar,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "27"),
	card.WithBonus(card.Bonus.Aember),
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
