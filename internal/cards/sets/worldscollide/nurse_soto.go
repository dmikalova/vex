//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// NurseSoto
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Human
//
//	Deploy. (This creature can enter play anywhere in your battleline.)
//	Play/Fight/Reap: Heal 3 damage from each of Nurse Soto's neighbors.
var NurseSoto = card.New(
	"Nurse Soto",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 315),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
