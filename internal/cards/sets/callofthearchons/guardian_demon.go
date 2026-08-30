//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// GuardianDemon
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Demon
//
//	Play/Fight/Reap: Heal up to 2 damage from a creature. Deal that amount of damage to another creature.
var GuardianDemon = card.New(
	"Guardian Demon",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 88),
	card.WithPower(4),
	card.WithTraits("Demon"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
