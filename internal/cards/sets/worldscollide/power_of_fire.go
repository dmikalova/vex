//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// PowerOfFire
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Sacrifice a friendly creature. If you do, each player loses A equal to half that creature's power (rounding down the loss). Gain 1 chain.
var PowerOfFire = card.New(
	"Power of Fire",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "26"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
