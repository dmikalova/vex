package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Combat Pheromones
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Item
//
//	Versatile.
//	Action: Destroy Combat Pheromones. Use 2 other Mars cards, one at a time.
var CombatPheromones = card.New(
	"Combat Pheromones",
	card.House.Mars,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 180),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Item),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.Action, card.Sentences{Effects: []card.Effect{
			card.Destroy{Target: card.Target.This},
			card.Use{
				Max:    2,
				Target: card.Target.EachFriendlyCardInPlay.OfHouse(card.House.Mars).Other(),
			},
		}}),
)
