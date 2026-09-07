//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// GronNineToes
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Giant
//
//	Gron Nine-Toes gets +4 power while it is damaged. (Gron Nine-Toes gets the power bonus only if he survives the damage.)
var GronNineToes = card.New(
	"Gron Nine-Toes",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, 9),
	card.WithPower(5),
	card.WithTraits(card.Traits.Giant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
