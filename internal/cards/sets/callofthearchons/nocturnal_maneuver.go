package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Nocturnal Maneuver
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Exhaust up to 3 creatures.
var NocturnalManeuver = card.New(
	"Nocturnal Maneuver",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 330),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.ExhaustCreatures{
			Max:    3,
			Target: card.Target.EachCreature,
		}),
)
