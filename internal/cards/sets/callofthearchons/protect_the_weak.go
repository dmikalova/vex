package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Protect the Weak
//
//	House:  Sanctum
//	Type:   Upgrade
//	Rarity: Common
//	Æmber:  1
//
//	This Creature gains +1 armor and taunt.
var ProtectTheWeak = card.New(
	"Protect the Weak",
	card.House.Sanctum,
	card.Type.Upgrade,
	card.Rarity.Common,
	card.Provenance(card.CotA, "265"),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		ArmorBonus: 1,
		Keywords:   card.Keywords(card.Keyword.Taunt),
	}),
)
