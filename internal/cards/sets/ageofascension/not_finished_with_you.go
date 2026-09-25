package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Not Finished with You
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Shuffle any number of creatures from your discard pile into your deck.
var NotFinishedWithYou = set.New(
	"Not Finished with You",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, "63"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.ShuffleIntoDeck{
			Player: card.Controller, From: []card.Zone{card.Discard},
			Selection: card.Chosen{
				Type:     card.Type.Creature,
				Optional: true,
			},
			Quantity: card.AnyNumber{},
		}),
)
