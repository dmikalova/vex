package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Not Finished with You
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Shuffle any number of Creatures from your discard pile into your deck.
var NotFinishedWithYou = card.New(
	"Not Finished with You",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, "63"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.ShuffleFromDiscard{
			Selection: card.Chosen{Type: card.Type.Creature, Optional: true},
			AnyNumber: true,
		}),
)
