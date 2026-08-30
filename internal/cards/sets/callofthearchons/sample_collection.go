//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// SampleCollection
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Put an enemy creature into your archives for each key your opponent has forged. If any of these creatures leave your archives, they are put into their owner's hand instead.
var SampleCollection = card.New(
	"Sample Collection",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 175),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
