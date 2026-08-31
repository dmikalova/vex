package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Feeding Pit
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Action: Discard a creature from your hand -> gain 1 Æmber.
var FeedingPit = card.New(
	"Feeding Pit",
	card.House.Mars,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 184),
	card.WithTraits("Location"),
	card.WithAbility(card.Trigger.Action, card.Then{
		First: card.DiscardFromHand{
			Count:         1,
			CreaturesOnly: true,
		},
		Result: card.GainAember{
			Player: card.Controller,
			Amount: 1,
		},
	}),
)
