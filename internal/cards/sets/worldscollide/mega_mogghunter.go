//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// MegaMogghunter
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: FIXED
//	Power:  8
//	Traits: Giant
//
//	Fight: Deal 2D to a flank creature.
var MegaMogghunter = card.New(
	"Mega Mogghunter",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.FIXED,
	card.Provenance(card.WC, 59),
	card.WithPower(8),
	card.WithTraits(card.Traits.Giant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
