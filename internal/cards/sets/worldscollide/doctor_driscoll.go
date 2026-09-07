//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// DoctorDriscoll
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Human • Scientist
//
//	Elusive.
//	Action: Heal 2 damage from a creature. Gain 1A for each damage healed this way.
var DoctorDriscoll = card.New(
	"Doctor Driscoll",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, 329),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human, card.Traits.Scientist),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
