package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Key Charge
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Lose 1 Æmber -> forge a key at current cost.
var KeyCharge = card.New(
	"Key Charge",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 325),
	card.WithAbility(card.Trigger.Play, card.Then{
		First: card.LoseAember{
			Player: card.Controller,
			Amount: 1,
		},
		Result: card.ForgeKey{},
	}),
)
