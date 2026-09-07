//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// GoodOfTheMany
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Destroy each creature that does not share a trait with another creature in its controller's battleline.
var GoodOfTheMany = card.New(
	"Good of the Many",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, "220"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
