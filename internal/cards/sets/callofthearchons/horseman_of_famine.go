package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Horseman of Famine
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Connected
//	Power:  5
//	Traits: Horseman • Spirit
//
//	Play/Fight/Reap: Destroy the least powerful creature.
var HorsemanOfFamine = set.New(
	"Horseman of Famine",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Connected,
	card.Provenance(card.CotA, "247"),
	card.InCluster(horsemenCluster),
	card.WithPower(5),
	card.WithTraits(card.Traits.Horseman, card.Traits.Spirit),
	card.WithAbility(card.Trigger.PlayFightReap, card.Destroy{
		Target: card.Target.EachCreature.Refine(card.LeastPowerful),
	}),
)
