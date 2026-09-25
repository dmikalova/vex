package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Shoulder Armor
//
//	House:  Sanctum
//	Type:   Upgrade
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	While this creature is on a flank, it gains +2 power and +2 armor.
var ShoulderArmor = set.New(
	"Shoulder Armor",
	card.House.Sanctum,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "266"),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{
		PowerBonus:   2,
		ArmorBonus:   2,
		WhileOnFlank: true,
	}),
)
