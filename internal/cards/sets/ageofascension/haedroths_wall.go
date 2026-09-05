package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Haedroth's Wall
//
//	House:  Sanctum
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Location
//
//	Each friendly flank creature gains +2 power.
var HaedrothsWall = card.New(
	"Haedroth's Wall",
	card.House.Sanctum,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 236),
	card.WithTraits(card.Traits.Location),
	card.WithConstant(card.ConstantAbility{
		PowerBonus: 2,
		Target:     card.Target.EachFriendlyCreature.OnFlank(),
	}),
)
