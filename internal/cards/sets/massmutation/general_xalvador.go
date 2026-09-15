package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// General Xalvador
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Armor:  2
//	Traits: Human • Knight
//
//	Enhance Capture Capture.
var GeneralXalvador = set.New(
	"General Xalvador",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "135"),
	card.WithPower(4),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Human, card.Traits.Knight),
	card.WithEnhance(card.Bonus.Capture, card.Bonus.Capture),
)
