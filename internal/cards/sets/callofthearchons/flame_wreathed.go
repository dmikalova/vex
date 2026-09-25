package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Flame-Wreathed
//
//	House:  Dis
//	Type:   Upgrade
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	This creature gains +2 power and +2 hazardous.
var FlameWreathed = set.New(
	"Flame-Wreathed",
	card.House.Dis,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "106"),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{
		PowerBonus:     2,
		HazardousBonus: 2,
	}),
)
