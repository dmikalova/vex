//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// DinosaursBane
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Variant
//	Æmber:  1
//
//	Play: Destroy a Dinosaur creature.
var DinosaursBane = card.New(
	"Dinosaurs' Bane",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, 125),
	card.WithAemberBonus(1),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
