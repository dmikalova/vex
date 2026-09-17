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
var FeedingPit = set.New(
	"Feeding Pit",
	card.House.Mars,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "184"),
	card.WithTraits(card.Traits.Location),
	card.WithAbility(card.Trigger.Action, card.Then{
		First: card.DiscardCard{
			Player:    card.Controller,
			Zones:     []card.Zone{card.Hand},
			Selection: card.Chosen{Type: card.Type.Creature},
			Amount:    1,
		},
		Result: card.GainAember{
			Player: card.Controller,
			Amount: 1,
		},
	}),
)
