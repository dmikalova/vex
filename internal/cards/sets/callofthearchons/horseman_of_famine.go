package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Horseman of Famine
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Connected
//	Power:  5
//	Traits: Horseman • Spirit
//
//	Play/Fight/Reap: Destroy the least powerful creature.
var HorsemanOfFamine = card.New(
	"Horseman of Famine",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Connected,
	card.Provenance(card.CotA, 247),
	card.WithPower(5),
	card.WithTraits("Horseman", "Spirit"),
	card.WithPlayFightReap(card.Destroy{
		Target: card.Target.EachCreature.Selector(card.LeastPowerful),
	}),
)
