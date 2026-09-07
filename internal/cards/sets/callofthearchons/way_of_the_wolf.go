package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Way of the Wolf
//
//	House:  Untamed
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	This creature gains skirmish.
var WayOfTheWolf = card.New(
	"Way of the Wolf",
	card.House.Untamed,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "370"),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{Keywords: card.Keywords(card.Keyword.Skirmish)}),
)
