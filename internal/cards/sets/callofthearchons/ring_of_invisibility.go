package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Ring of Invisibility
//
//	House:  Shadows
//	Type:   Upgrade
//	Rarity: Rare
//	Bonus:  Æmber
//
//	This creature gains elusive and skirmish.
var RingOfInvisibility = set.New(
	"Ring of Invisibility",
	card.House.Shadows,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "317"),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(
		card.StaticModifier{Keywords: card.Keywords(card.Keyword.Elusive, card.Keyword.Skirmish)},
	),
)
