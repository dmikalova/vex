package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Song of Spring
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Shuffle any number of friendly Untamed creatures from your hand, discard pile, or battleline into your deck.
var SongOfSpring = card.New(
	"Song of Spring",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, "332"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.ShuffleChosenCreaturesFromZones{
			House: card.House.Self,
		}),
)
