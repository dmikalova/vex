package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Mind Barb
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Your opponent discards a random card from their hand.
var MindBarb = card.New(
	"Mind Barb",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 67),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.DiscardRandomFromHand{Player: card.Opponent}),
)
