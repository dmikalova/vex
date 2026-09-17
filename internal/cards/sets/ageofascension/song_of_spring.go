package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Song of Spring
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Shuffle any number of friendly Untamed creatures from your hand, discard pile, or battleline into your deck.
var SongOfSpring = set.New(
	"Song of Spring",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, "332"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.ShuffleChosenCreaturesFromZones{
			House: card.Houses.Named(card.House.Self),
		}),
)
