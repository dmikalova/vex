package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Lieutenant Gorvenal
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Armor:  1
//	Traits: Spirit • Knight
//
//	After a friendly creature is used to fight, Lieutenant Gorvenal captures 1 Æmber from your opponent.
var LieutenantGorvenal = set.New(
	"Lieutenant Gorvenal",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "137"),
	card.WithPower(4),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Spirit, card.Traits.Knight),
	card.WithAbility(
		card.Trigger.AfterCreatureFights, card.Conditional{
			Cond: card.ItIsFriendly{},
			Then: card.CaptureAember{
				Amount: 1,
				Target: card.Target.This,
				Source: card.Opponent,
			},
		}),
)
