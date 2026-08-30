//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// NiffleQueen
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  6
//	Traits: Beast • Niffle
//
//	Each other friendly Beast creature gets +1 power.
//	Each other friendly Niffle creature gets +1 power.
var NiffleQueen = card.New(
	"Niffle Queen",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 364),
	card.WithPower(6),
	card.WithTraits("Beast", "Niffle"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
