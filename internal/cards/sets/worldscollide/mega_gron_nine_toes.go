//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// MegaGronNineToes
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
//	Mega Gron Nine-Toes gets +4 power while it is damaged. (Mega Gron Nine-Toes gets the power bonus only if he survives the damage.)
var MegaGronNineToes = card.New(
	"Mega Gron Nine-Toes",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.FIXED,
	card.Provenance(card.WC, 58),
	card.WithPower(7),
	card.WithTraits(card.Traits.Giant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
