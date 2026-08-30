//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// RelentlessAssault
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Ready and fight with up to 3 different friendly creatures, one at a time.
var RelentlessAssault = card.New(
	"Relentless Assault",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 13),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
