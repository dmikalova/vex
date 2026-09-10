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
//	Fight/Reap: For the remainder of the turn, you may play or use a Brobnar card.
var BrobnarAmbassador = card.New(
	"Brobnar Ambassador",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Special,
	card.Provenance(card.AoA, "229"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Human),
	card.WithKeywords(card.Keyword.Elusive),
	// TODO: planned rework of the Ambassador cycle.
	card.WithAbility(card.Trigger.FightReap,
		card.MayActFriendlyHouse{House: card.House.Brobnar, Grant: card.GrantPlay | card.GrantUse},
	),
)
