package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Bordan the Redeemed
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Human • Thief
//
//	Elusive.
//	Action: Bordan the Redeemed captures 2 Æmber from your opponent.
var BordanTheRedeemed = card.New(
	"Bordan the Redeemed",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 215),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human, card.Traits.Thief),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Action, card.CaptureAember{
			Amount: 2,
			Target: card.Target.This,
			Source: card.Opponent,
		}),
)
