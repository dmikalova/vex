package anomalyexpansion

import "github.com/dmikalova/vactrol/internal/card"

// Orb of Wonder
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Special
//	Traits: Item
//
//	Versatile.
//	Action: Destroy Orb of Wonder -> search your deck for a card and put it into your hand. Shuffle your deck.
var OrbOfWonder = set.New(
	"Orb of Wonder",
	card.House.Brobnar,
	card.Type.Artifact,
	// Rarity relabelled from FIXED to Special.
	card.Rarity.Special,
	card.Provenance(card.WC, "A06"),
	card.WithTraits(card.Traits.Item),
	card.WithKeywords(card.Keyword.Versatile),
	card.WithAbility(
		card.Trigger.Action, card.Then{
			First: card.Destroy{Target: card.Target.This},
			Result: card.Sentences{
				Effects: []card.Effect{
					card.SearchDeck{},
					card.Shuffle{},
				},
			},
		}),
)
