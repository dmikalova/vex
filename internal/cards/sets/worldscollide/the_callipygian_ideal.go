//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// TheCallipygianIdeal
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Upgrade
//	Rarity: Uncommon
//
//	Play: Exalt this creature.
//	This creature gains, "You may spend A on this creature as if it were in your pool."
var TheCallipygianIdeal = card.New(
	"The Callipygian Ideal",
	card.House.Saurian,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 212),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
