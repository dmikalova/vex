package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Xenos Bloodshadow
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Human • Witch
//
//	Elusive, Poison, Skirmish, Hazardous 6.
var XenosBloodshadow = card.New(
	"Xenos Bloodshadow",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, 404),
	card.WithPower(4),
	card.WithTraits(card.Traits.Human, card.Traits.Witch),
	card.WithKeywords(card.Keyword.Elusive, card.Keyword.Poison, card.Keyword.Skirmish),
	card.WithHazardous(6),
)
