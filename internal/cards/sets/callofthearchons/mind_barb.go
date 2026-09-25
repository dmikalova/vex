package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Mind Barb
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Discard a card from your hand. Your opponent discards a random card from their hand.
var MindBarb = set.New(
	"Mind Barb",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, "67"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(card.Trigger.Play, card.Sequence{
		Effects: []card.Effect{
			card.DiscardCard{
				Player:    card.Controller,
				Zones:     []card.Zone{card.Hand},
				Selection: card.Chosen{},
			},
			card.DiscardCard{
				Player:    card.Opponent,
				Zones:     []card.Zone{card.Hand},
				Selection: card.Random{},
			},
		},
	}),
)
