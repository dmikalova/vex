package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Yurk
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Demon
//
//	Play: Discard a card from your hand.
var Yurk = card.New(
	"Yurk",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 68),
	card.WithPower(4),
	card.WithTraits(card.Traits.Demon),
	card.WithAbility(
		card.Trigger.Play, card.DiscardFromHand{Count: 1}),
)
