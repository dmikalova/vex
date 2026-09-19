package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Combat Pheromones
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Uncommon
//	Bonus:  Æmber
//	Traits: Item
//
//	Versatile.
//	Action: Destroy Combat Pheromones. Use 2 other Mars cards, one at a time.
var CombatPheromones = set.New(
	"Combat Pheromones",
	card.House.Mars,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "180"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Item),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.Action, card.Sequence{Effects: []card.Effect{
			card.Destroy{Target: card.Target.This},
			card.Use{
				Max: 2,
				Target: card.Target.EachFriendlyCardInPlay.House(card.Houses.Named(card.House.Self)).
					Other(),
			},
		}}),
)
