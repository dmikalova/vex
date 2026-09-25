package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Virtuous Works
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber Æmber Æmber
var VirtuousWorks = set.New(
	"Virtuous Works",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "230"),
	card.WithBonus(card.Bonus.Aember, card.Bonus.Aember, card.Bonus.Aember),
)
