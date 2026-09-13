package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Mantle of the Zealot
//
//	House:  Sanctum
//	Type:   Upgrade
//	Rarity: Rare
//
//	This Creature gains versatile.
var MantleOfTheZealot = card.New(
	"Mantle of the Zealot",
	card.House.Sanctum,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "264"),
	card.WithStatic(card.StaticModifier{Keywords: card.Keywords(card.Keyword.Versatile)}),
)
