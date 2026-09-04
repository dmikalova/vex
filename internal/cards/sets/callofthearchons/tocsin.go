package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Tocsin
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Demon
//
//	Reap: Your opponent discards a random card from their hand.
var Tocsin = card.New(
	"Tocsin",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 102),
	card.WithPower(3),
	card.WithTraits(card.Traits.Demon),
	card.WithAbility(
		card.Trigger.Reap, card.DiscardRandomFromHand{Player: card.Opponent}),
)
