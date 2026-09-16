package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Titan Engineer
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  6
//	Traits: Cyborg • Scientist
//
//	While Titan Engineer is not on a flank, each player's keys cost +1 Æmber.
var TitanEngineer = set.New(
	"Titan Engineer",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "081"),
	card.WithPower(6),
	card.WithTraits(card.Traits.Cyborg, card.Traits.Scientist),
	card.WithKeyCost(card.KeyCostChange(card.EachPlayer, 1).WhileOffFlank()),
)
