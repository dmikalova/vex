package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Plasma Nozzle
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This Creature gains +2 assault and +2 splash-attack.
var PlasmaNozzle = set.New(
	"Plasma Nozzle",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, "336"),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		AssaultBonus:      2,
		SplashAttackBonus: 2,
	}),
)
