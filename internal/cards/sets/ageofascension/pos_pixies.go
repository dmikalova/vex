package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Po's Pixies
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Faerie
//
//	Elusive.
//	Æmber stolen or captured from your pool is taken from the common supply instead.
var PosPixies = card.New(
	"Po's Pixies",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 362),
	card.WithPower(1),
	card.WithTraits(card.Traits.Faerie),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithTheftFromCommonSupply(),
)
