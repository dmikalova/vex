//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// LateralShift
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Special
//
//	Play: Look at your opponent's hand. Play a card from that hand as if it were yours.
var LateralShift = card.New(
	"Lateral Shift",
	card.House.Brobnar,
	card.Type.Tactic,
	// TODO(variant): rarity relabelled from FIXED to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.WC, "A03"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
