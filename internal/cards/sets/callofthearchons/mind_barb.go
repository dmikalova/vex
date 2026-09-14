package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Mind Barb
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Discard a card from your hand. Your opponent discards a random card from their hand.
var MindBarb = set.New(
	"Mind Barb",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, "67"),
	card.WithAemberBonus(1),
	card.WithAbility(card.Trigger.Play, card.Sentences{
		Effects: []card.Effect{
			card.DiscardCard{
				Player:    card.Controller,
				Zone:      card.Hand,
				Selection: card.Chosen{},
				Amount:    1,
			},
			card.DiscardCard{Player: card.Opponent, Zone: card.Hand, Selection: card.Random{}},
		},
	}),
)
