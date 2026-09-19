package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Seismo-entangler
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Item
//
//	Action: Choose a house. Your opponent cannot use creatures of the chosen house to reap during their next turn.
var SeismoEntangler = set.New(
	"Seismo-entangler",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "137"),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.Action, card.ChooseHouseThen{
			Then: card.Restrict{
				Player:   card.Opponent,
				Action:   card.Restricted.Reaping,
				House:    card.TheChosenHouse,
				Duration: card.Duration.OpponentNextTurn,
			},
		}),
)
