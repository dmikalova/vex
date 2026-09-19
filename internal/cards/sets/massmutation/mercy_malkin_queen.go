package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Mercy, Malkin Queen
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Human • Witch
//
//	Skirmish.
//	After a creature enters play, if it is a friendly Cat creature, ward it.
//	Fight: Ready a friendly Beast creature.
var MercyMalkinQueen = set.New(
	"Mercy, Malkin Queen",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "403"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human, card.Traits.Witch),
	card.WithKeywords(card.Keyword.Skirmish),
	card.WithAbility(
		card.Trigger.AfterCreatureEnters, card.Conditional{
			Cond: card.And{Conditions: []card.Condition{
				card.ItIsFriendly{},
				card.ItIsOfTrait{Trait: card.Traits.Cat},
			}},
			Then: card.Ward{Target: card.Target.Triggering},
		}),
	card.WithAbility(
		card.Trigger.Fight, card.ReadyCreatures{
			Target: card.Target.EachFriendlyCreature.WithTrait(card.Traits.Beast),
		}),
)
