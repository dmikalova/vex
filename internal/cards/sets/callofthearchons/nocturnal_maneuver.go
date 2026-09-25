package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Nocturnal Maneuver
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Exhaust up to 3 creatures.
var NocturnalManeuver = set.New(
	"Nocturnal Maneuver",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, "330"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.ExhaustCreatures{
			Max:    3,
			Target: card.Target.EachCreature,
		}),
)
