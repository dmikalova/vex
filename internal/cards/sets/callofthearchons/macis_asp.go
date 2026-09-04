package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Macis Asp
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Beast
//
//	Skirmish, Poison.
var MacisAsp = card.New(
	"Macis Asp",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 301),
	card.WithPower(3),
	card.WithTraits(card.Traits.Beast),
	card.WithKeywords(card.Keyword.Skirmish, card.Keyword.Poison),
)
