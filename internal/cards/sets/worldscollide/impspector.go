package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Impspector
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Imp
//
//	Destroyed: Purge a random card from your opponent's hand.
var Impspector = set.New(
	"Impspector",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "77"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Imp),
	card.WithAbility(
		card.Trigger.Destroyed, card.PurgeFromHand{
			Player:    card.Opponent,
			Selection: card.Random{},
		}),
)
