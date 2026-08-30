package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Blinding Light
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Choose a house - stun each creature of the chosen house.
var BlindingLight = card.New(
	"Blinding Light",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 213),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.ChooseHouseThen{
			Then: card.Stun{Target: card.Target.EachCreature.OfChosenHouse()},
		}),
)
