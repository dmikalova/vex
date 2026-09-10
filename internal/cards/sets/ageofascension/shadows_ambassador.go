package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Shadows Ambassador
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Special
//	Power:  1
//	Traits: Human
//
//	Elusive.
//	Fight/Reap: For the remainder of the turn, you may play or use a Shadows card.
var ShadowsAmbassador = card.New(
	"Shadows Ambassador",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Special,
	card.Provenance(card.AoA, "243"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Human),
	card.WithKeywords(card.Keyword.Elusive),
	// TODO: planned rework of the Ambassador cycle.
	card.WithAbility(card.Trigger.FightOrReap,
		card.MayActFriendlyHouse{House: card.House.Shadows, Grant: card.GrantPlay | card.GrantUse},
	),
)
