//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// GloriousFew
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: For each creature your opponent controls in excess of you, gain 1 Aember.
var GloriousFew = card.New(
	"Glorious Few",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 218),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
