package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Nexus
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Cyborg • Thief
//
//	Elusive.
//	Reap: Use an enemy artifact.
var Nexus = card.New(
	"Nexus",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 305),
	card.WithPower(3),
	card.WithTraits("Cyborg", "Thief"),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Reap, card.Use{
			Max:    1,
			Target: card.Target.EachEnemyArtifact,
		}),
)
