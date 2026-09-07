//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// TachyonPulse
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Destroy each artifact. Exhaust each creature with an upgrade.
var TachyonPulse = card.New(
	"Tachyon Pulse",
	card.House.Staralliance,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, 340),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
