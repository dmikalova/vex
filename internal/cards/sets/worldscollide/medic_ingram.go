//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// MedicIngram
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Human
//
//	Play/Fight/Reap: You may heal 3 damage from a creature and ward it.
var MedicIngram = card.New(
	"Medic Ingram",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, 301),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
