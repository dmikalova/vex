package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Raiding Knight
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Human • Knight
//
//	Play: Raiding Knight captures 1 Æmber from your opponent.
var RaidingKnight = card.New(
	"Raiding Knight",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 255),
	card.WithPower(4),
	card.WithTraits("Human", "Knight"),
	card.WithAbility(
		card.Trigger.Play, card.CaptureAember{
			Amount: 1,
			Target: card.Target.This,
			Source: card.Opponent,
		}),
)
