package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Charette
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Demon
//
//	Play: Charette captures 3 Æmber from your opponent.
var Charette = card.New(
	"Charette",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 81),
	card.WithPower(4),
	card.WithTraits(card.Traits.Demon),
	card.WithAbility(
		card.Trigger.Play, card.CaptureAember{
			Amount: 3,
			Target: card.Target.This,
			Source: card.Opponent,
		}),
)
