package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// World Tree
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Action: Put a creature from your discard pile on top of your deck.
var WorldTree = set.New(
	"World Tree",
	card.House.Untamed,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "344"),
	card.WithTraits(card.Traits.Location),
	card.WithAbility(
		card.Trigger.Action, card.PutCard{Zones: []card.Zone{card.Discard},
			Selection:   card.Chosen{Type: card.Type.Creature},
			Destination: card.To.TopOfDeck,
		}),
)
