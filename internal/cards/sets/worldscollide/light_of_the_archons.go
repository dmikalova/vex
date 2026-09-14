package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Light of the Archons
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Common
//	Æmber:  1
//
//	This Creature gains +1 power and +1 armor for each Upgrade attached to it.
var LightOfTheArchons = set.New(
	"Light of the Archons",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Common,
	card.Provenance(card.WC, "300"),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		PowerBonus: 1,
		ArmorBonus: 1,
		Per:        card.UpgradesOnIt,
	}),
)
