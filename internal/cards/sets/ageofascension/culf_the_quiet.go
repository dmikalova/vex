package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Culf the Quiet
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  6
//	Traits: Giant
//
//	Elusive.
var CulfTheQuiet = card.New(
	"Culf the Quiet",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 20),
	card.WithPower(6),
	card.WithTraits(card.Traits.Giant),
	card.WithKeywords(card.Keyword.Elusive),
)
