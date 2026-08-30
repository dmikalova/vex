package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Snufflegator
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Beast
//
//	Skirmish.
var Snufflegator = card.New(
	"Snufflegator",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 358),
	card.WithPower(4),
	card.WithTraits("Beast"),
	card.WithKeywords(card.Keyword.Skirmish),
)
