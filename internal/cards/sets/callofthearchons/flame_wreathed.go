package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Flame-Wreathed
//
//	House:  Dis
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	This Creature gains +2 power and +2 hazardous.
var FlameWreathed = set.New(
	"Flame-Wreathed",
	card.House.Dis,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "106"),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		PowerBonus:     2,
		HazardousBonus: 2,
	}),
)
