package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Sanctum Guardian
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Traits: Knight • Spirit
//
//	Taunt.
//	Fight/Reap: Swap this creature with another friendly creature in your battleline.
var SanctumGuardian = card.New(
	"Sanctum Guardian",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 256),
	card.WithPower(6),
	card.WithTraits("Knight", "Spirit"),
	card.WithKeywords(card.Keyword.Taunt),
	card.WithFightOrReap(card.Swap{With: card.Target.OtherFriendlyCreature}),
)
