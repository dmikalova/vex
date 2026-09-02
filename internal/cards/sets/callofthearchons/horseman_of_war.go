package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Horseman of War
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Connected
//	Power:  5
//	Traits: Horseman • Spirit
//
//	Play: For the remainder of the turn, each friendly creature may fight.
var HorsemanOfWar = card.New(
	"Horseman of War",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Connected,
	card.Provenance(card.CotA, 249),
	card.WithPower(5),
	card.WithTraits("Horseman", "Spirit"),
	card.WithAbility(card.Trigger.Play, card.GrantFightAnyHouse{}),
)
