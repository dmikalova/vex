package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Blinding Light
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Choose a house. Stun each creature of the chosen house.
var BlindingLight = set.New(
	"Blinding Light",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, "213"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.ChooseHouseThen{
			Then: card.Stun{Target: card.Target.EachCreature.House(card.Houses.Chosen)},
		}),
)
