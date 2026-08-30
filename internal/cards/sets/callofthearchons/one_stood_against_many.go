//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// OneStoodAgainstMany
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Ready and fight with a friendly creature 3 times, each time against a different enemy creature. Resolve these fights one at a time.
var OneStoodAgainstMany = card.New(
	"One Stood Against Many",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 223),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
