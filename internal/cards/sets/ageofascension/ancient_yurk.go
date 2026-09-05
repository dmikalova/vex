package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Ancient Yurk
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Traits: Demon
//
//	Play: Discard 3 cards from your hand.
var AncientYurk = card.New(
	"Ancient Yurk",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 88),
	card.WithPower(6),
	card.WithTraits(card.Traits.Demon),
	card.WithAbility(
		card.Trigger.Play, card.DiscardFromHand{Count: 3}),
)
