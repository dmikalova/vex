package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Song of Spring
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Shuffle any number of friendly Untamed creatures from your hand, your discard pile, or play into your deck.
var SongOfSpring = set.New(
	"Song of Spring",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, "332"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.ShuffleIntoDeck{
			Player: card.Controller, From: []card.Zone{card.Hand, card.Discard, card.InPlay},
			Selection: card.Chosen{
				House:    card.Houses.Named(card.House.Self),
				Type:     card.Type.Creature,
				Optional: true,
			},
			Quantity: card.AnyNumber{},
		}),
)
