package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Lights Out
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Put up to 2 creatures into their owners' hands.
var LightsOut = card.New(
	"Lights Out",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 274),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.PutUpTo{
			Max:         2,
			Target:      card.Target.EachEnemyCreature,
			Destination: card.To.Hand,
		}),
)
