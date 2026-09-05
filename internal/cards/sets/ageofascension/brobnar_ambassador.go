package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Brobnar Ambassador
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Special
//	Power:  1
//	Traits: Human
//
//	Elusive.
//	Fight/Reap: You may play or use a Brobnar card this turn.
var BrobnarAmbassador = card.New(
	"Brobnar Ambassador",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Special,
	card.Provenance(card.AoA, 229),
	card.WithPower(1),
	card.WithTraits(card.Traits.Human),
	card.WithKeywords(card.Keyword.Elusive),
	// TODO: planned rework of the Ambassador cycle.
	card.WithFightOrReap(card.MayPlayOrUseFriendlyHouse{House: card.House.Brobnar}),
)
