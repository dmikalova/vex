//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// OathOfPoverty
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Destroy each of your artifacts. Gain 2 Aember for each artifact destroyed this way.
var OathOfPoverty = card.New(
	"Oath of Poverty",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 222),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
