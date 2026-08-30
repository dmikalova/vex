package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Doorstep to Heaven
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Each player with 6 Æmber or more loses all but 5 Æmber.
var DoorstepToHeaven = card.New(
	"Doorstep to Heaven",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 217),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.LoseAember{
			Player: card.EachPlayer,
			By:     card.AllBut(5),
		}),
)
