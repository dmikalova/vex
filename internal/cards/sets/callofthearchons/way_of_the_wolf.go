package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Way of the Wolf
//
//	House:  Untamed
//	Type:   Upgrade
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	This creature gains skirmish.
var WayOfTheWolf = set.New(
	"Way of the Wolf",
	card.House.Untamed,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "370"),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{Keywords: card.Keywords(card.Keyword.Skirmish)}),
)
