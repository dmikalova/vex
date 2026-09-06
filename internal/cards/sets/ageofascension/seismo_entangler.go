package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Seismo-entangler
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Item
//
//	Action: Choose a house - during your opponent's next turn, creatures of the chosen house cannot be used to reap.
var SeismoEntangler = card.New(
	"Seismo-entangler",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 137),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.Action, card.ChooseHouseThen{
			Then: card.ChosenHouseCannotReapNextTurn{Player: card.Opponent},
		}),
)
