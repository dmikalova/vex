package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Titan Mechanic
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  6
//	Traits: Cyborg • Scientist
//
//	While Titan Mechanic is on a flank, each player's keys cost -1 Æmber.
var TitanMechanic = card.New(
	"Titan Mechanic",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 154),
	card.WithPower(6),
	card.WithTraits("Cyborg", "Scientist"),
	card.WithKeyCost(card.KeyCostChange(card.EachPlayer, -1).WhileOnFlank()),
)
