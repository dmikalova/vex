//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// CrashMuldoon
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Human • Pilot
//
//	Deploy.
//	Crash Muldoon enters play ready.
//	Action: Use a neighboring non-Star Alliance creature.
var CrashMuldoon = card.New(
	"Crash Muldoon",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "327"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human, card.Traits.Pilot),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
