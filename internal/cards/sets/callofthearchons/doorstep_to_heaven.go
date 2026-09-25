package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Doorstep to Heaven
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Each player with 6 Æmber or more loses all but 5 Æmber.
var DoorstepToHeaven = set.New(
	"Doorstep to Heaven",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "217"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.LoseAember{
			Player: card.EachPlayer,
			By:     card.AllBut(5),
		}),
)
