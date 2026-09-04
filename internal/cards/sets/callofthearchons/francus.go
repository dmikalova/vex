package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Francus
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  6
//	Armor:  1
//	Traits: Knight • Spirit
//
//	After a creature is destroyed fighting Francus, Francus captures 1 Æmber from your opponent.
var Francus = card.New(
	"Francus",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 243),
	card.WithPower(6),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Knight, card.Traits.Spirit),
	card.WithAbility(
		card.Trigger.AfterDestroyedFighting, card.CaptureAember{
			Amount: 1,
			Target: card.Target.This,
			Source: card.Opponent,
		}),
)
