//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// EeOnTheFringes
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  1
//	Traits: Imp
//
//	Elusive.
//	During your turn, after you discard a Dis card from your hand, you may purge a Dis card from a discard pile. If you do, steal 1A.
var EeOnTheFringes = card.New(
	"E'e on the Fringes",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "88"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Imp),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
