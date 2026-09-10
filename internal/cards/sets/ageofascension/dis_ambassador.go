package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Dis Ambassador
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Special
//	Power:  1
//	Traits: Human
//
//	Elusive.
//	Fight/Reap: For the remainder of the turn, you may play or use a Dis card.
var DisAmbassador = card.New(
	"Dis Ambassador",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Special,
	card.Provenance(card.AoA, "230"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Human),
	card.WithKeywords(card.Keyword.Elusive),
	// TODO: planned rework of the Ambassador cycle.
	card.WithAbility(card.Trigger.FightOrReap,
		card.MayActFriendlyHouse{House: card.House.Dis, Grant: card.GrantPlay | card.GrantUse},
	),
)

// TODO: should not be special
