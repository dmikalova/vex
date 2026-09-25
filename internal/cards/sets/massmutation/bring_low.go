package massmutation

import "github.com/dmikalova/vex/internal/card"

// Bring Low
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Capture all but 5 Æmber from your opponent, distributed among any number of friendly creatures.
//	Enhance Capture.
var BringLow = set.New(
	"Bring Low",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "147"),
	card.WithBonus(card.Bonus.Aember),
	card.WithEnhance(card.Bonus.Capture),
	card.WithAbility(
		card.Trigger.Play, card.DistributeCapture{
			By:     card.AllBut(5),
			Source: card.Opponent,
		}),
)
