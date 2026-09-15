package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Squire Alys
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Armor:  2
//	Traits: Human
//
//	Play: Squire Alys captures 2 Æmber from your opponent.
var SquireAlys = set.New(
	"Squire Alys",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "140"),
	card.WithPower(2),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Human),
	card.WithAbility(
		card.Trigger.Play, card.CaptureAember{
			Amount: 2,
			Target: card.Target.This,
			Source: card.Opponent,
		}),
)
