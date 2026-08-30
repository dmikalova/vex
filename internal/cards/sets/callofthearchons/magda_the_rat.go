//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// MagdaTheRat
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Elf • Thief
//
//	Elusive. (The first time this creature is attacked each turn, no damage is dealt.)
//	Play: Steal 2 Aember.
//	Leaves Play: Your opponent steals 2 Aember.
var MagdaTheRat = card.New(
	"Magda the Rat",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 303),
	card.WithPower(4),
	card.WithTraits("Elf", "Thief"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
