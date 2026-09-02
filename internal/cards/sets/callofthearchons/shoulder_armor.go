package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Shoulder Armor
//
//	House:  Sanctum
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	While this creature is on a flank, it gains +2 power and +2 armor.
var ShoulderArmor = card.New(
	"Shoulder Armor",
	card.House.Sanctum,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 266),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		PowerBonus:   2,
		ArmorBonus:   2,
		WhileOnFlank: true,
	}),
)
