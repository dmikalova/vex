package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Wrath
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Special
//	Power:  3
//	Armor:  3
//	Traits: Demon • Sin
//
//	Taunt, Poison, Skirmish.
//	Fight: For each friendly Sin creature, enrage an enemy creature.
var Wrath = set.New(
	"Wrath",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Special,
	card.Provenance(card.MM, "064"),
	card.InCluster(sinsCluster),
	card.OneCopyPerDeck(),
	card.WithPower(3),
	card.WithArmor(3),
	card.WithTraits(card.Traits.Demon, card.Traits.Sin),
	card.WithKeywords(card.Keyword.Taunt, card.Keyword.Poison, card.Keyword.Skirmish),
	card.WithAbility(card.Trigger.Fight, card.ForEach{
		Times: card.CardsInPlay{
			Player: card.Controller,
			Type:   card.Type.Creature,
			Trait:  card.Traits.Sin,
		},
		Do: card.Enrage{Target: card.Target.EnemyCreature},
	}),
)
