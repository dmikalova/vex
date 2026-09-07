//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// MegaGangerChieftain
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: FIXED
//	Power:  7
//	Traits: Giant
//
//	Play: You may ready and fight with a neighboring creature.
var MegaGangerChieftain = card.New(
	"Mega Ganger Chieftain",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.FIXED,
	card.Provenance(card.WC, 56),
	card.WithPower(7),
	card.WithTraits(card.Traits.Giant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
