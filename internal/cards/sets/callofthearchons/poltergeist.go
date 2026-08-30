//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Poltergeist
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Use an artifact controlled by any player as if it were yours. Destroy that artifact.
var Poltergeist = card.New(
	"Poltergeist",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 69),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
