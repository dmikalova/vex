package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Orb of Wonder
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Versatile.
//	Action: Destroy Orb of Wonder -> search your deck for a card and put it into your hand. Shuffle your deck.
var OrbOfWonder = set.New(
	"Orb of Wonder",
	card.House.Sanctum,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.MM, "173"),
	card.WithTraits(card.Traits.Item),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.Action, card.Then{
			First: card.Destroy{Target: card.Target.This},
			Result: card.Sequence{
				Effects: []card.Effect{
					card.Search{
						Sources: []card.Zone{card.Deck},
						Dest:    card.To.Hand,
					},
					card.Shuffle{},
				},
			},
		}),
)
