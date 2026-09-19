package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Fetchdrones
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: Discard the top 2 cards of your deck. For each Logos card discarded this way, a friendly creature captures 2 Æmber from your opponent.
var Fetchdrones = set.New(
	"Fetchdrones",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "144"),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.Action, card.Sequence{Effects: []card.Effect{
			card.DiscardTop{
				Player: card.Controller,
				Amount: 2,
			},
			card.ForEachDiscarded{
				House: card.Houses.Named(card.House.Self),
				Do: card.CaptureAember{
					Amount: 2,
					Target: card.Target.FriendlyCreature,
					Source: card.Opponent,
				},
			},
		}}),
)
