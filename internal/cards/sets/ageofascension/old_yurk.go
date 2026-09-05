package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Old Yurk
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Traits: Demon
//
//	Play: Discard 2 cards from your hand.
var OldYurk = card.New(
	"Old Yurk",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 77),
	card.WithPower(5),
	card.WithTraits(card.Traits.Demon),
	card.WithAbility(
		card.Trigger.Play, card.DiscardFromHand{Count: 2}),
)
