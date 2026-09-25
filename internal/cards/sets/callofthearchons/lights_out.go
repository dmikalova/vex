package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Lights Out
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Put up to 2 enemy creatures into their owners' hands.
var LightsOut = set.New(
	"Lights Out",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "274"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.PutChosen{
			Quantity:    card.UpTo{N: card.Fixed(2)},
			Target:      card.Target.EachEnemyCreature,
			Destination: card.To.Hand,
		}),
)
