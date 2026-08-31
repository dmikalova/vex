package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Ring of Invisibility
//
//	House:  Shadows
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This creature gains elusive and skirmish.
var RingOfInvisibility = card.New(
	"Ring of Invisibility",
	card.House.Shadows,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 317),
	card.WithAemberBonus(1),
	card.WithStatic(
		card.StaticModifier{Keywords: card.Keywords(card.Keyword.Elusive, card.Keyword.Skirmish)},
	),
)
