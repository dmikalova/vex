package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Blood of Titans
//
//	House:  Brobnar
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	This Creature gains +5 power.
var BloodOfTitans = set.New(
	"Blood of Titans",
	card.House.Brobnar,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "50"),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{PowerBonus: 5}),
)
