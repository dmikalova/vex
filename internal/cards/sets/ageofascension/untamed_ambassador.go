package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Untamed Ambassador
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Special
//	Power:  1
//	Traits: Human
//
//	Elusive.
//	Fight/Reap: You may play or use a Untamed card this turn.
var UntamedAmbassador = card.New(
	"Untamed Ambassador",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Special,
	card.Provenance(card.AoA, 247),
	card.WithPower(1),
	card.WithTraits(card.Traits.Human),
	card.WithKeywords(card.Keyword.Elusive),
	// TODO: planned rework of the Ambassador cycle.
	card.WithFightOrReap(card.MayPlayOrUseFriendlyHouse{House: card.House.Untamed}),
)
