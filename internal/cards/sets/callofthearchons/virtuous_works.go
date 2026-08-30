package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Virtuous Works
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  3
var VirtuousWorks = card.New(
	"Virtuous Works",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 230),
	card.WithAemberBonus(3),
)
