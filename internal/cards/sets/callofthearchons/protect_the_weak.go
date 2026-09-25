package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Protect the Weak
//
//	House:  Sanctum
//	Type:   Upgrade
//	Rarity: Common
//	Bonus:  Æmber
//
//	This creature gains +1 armor and taunt.
var ProtectTheWeak = set.New(
	"Protect the Weak",
	card.House.Sanctum,
	card.Type.Upgrade,
	card.Rarity.Common,
	card.Provenance(card.CotA, "265"),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{
		ArmorBonus: 1,
		Keywords:   card.Keywords(card.Keyword.Taunt),
	}),
)
