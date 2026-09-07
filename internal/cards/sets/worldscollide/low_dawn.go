//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// LowDawn
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: If there are 3 or more Untamed creatures in your discard pile, gain 2A. Shuffle each Untamed creature from your discard pile into your deck.
var LowDawn = card.New(
	"Low Dawn",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "377"),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
