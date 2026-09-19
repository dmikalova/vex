package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Book of Malefaction
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Item • Law
//
//	Versatile.
//	After Æmber is stolen from you, for each Æmber stolen, put a warrant counter on Book of Malefaction.
//	Action: Remove a warrant counter from Book of Malefaction -> purge a creature.
var BookOfMalefaction = set.New(
	"Book of Malefaction",
	card.House.Sanctum,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.MM, "159"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Item, card.Traits.Law),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.AfterAemberStolenFromYou, card.PlaceCounter{Amount: 1,
			Kind:   card.Counter.Warrant,
			Target: card.Target.This,
			Per:    card.AemberStolenThisEvent{},
		}),
	card.WithAbility(
		card.Trigger.Action, card.Then{
			First: card.RemoveCounters{
				Kind:   card.Counter.Warrant,
				Target: card.Target.This,
				Amount: 1,
			},
			Result: card.PurgeCreature{Target: card.Target.Creature},
		}),
)
