package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Desire
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Special
//	Power:  3
//	Traits: Demon • Sin
//
//	Each player's keys cost +4 Æmber.
//	Reap: Forge a key at current cost, reduced by 1 Æmber for each friendly Sin creature.
var Desire = set.New(
	"Desire",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Special,
	card.Provenance(card.MM, "053"),
	card.InCluster(sinsCluster),
	card.OneCopyPerDeck(),
	card.WithPower(3),
	card.WithTraits(card.Traits.Demon, card.Traits.Sin),
	card.WithKeyCost(card.KeyCostChange(card.EachPlayer, 4)),
	card.WithAbility(card.Trigger.Reap, card.ForgeKey{
		Discount: true,
		Keep:     true,
		ReducedBy: card.InPlay{
			Player: card.Controller,
			Type:   card.Type.Creature,
			Trait:  card.Traits.Sin,
		},
	}),
)
