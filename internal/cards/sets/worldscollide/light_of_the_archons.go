package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Light of the Archons
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Common
//	Bonus:  Æmber
//
//	This creature gains +1 power and +1 armor for each upgrade attached to it.
var LightOfTheArchons = set.New(
	"Light of the Archons",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Common,
	card.Provenance(card.WC, "300"),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{
		PowerBonus: 1,
		ArmorBonus: 1,
		Per:        card.UpgradesOnIt,
	}),
)
