//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Sneklifter
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Elf • Thief
//
//	Play: Take control of an enemy artifact. While under your control, if it does not belong to one of your three houses, it is considered to be of house Shadows.
var Sneklifter = card.New(
	"Sneklifter",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 313),
	card.WithPower(2),
	card.WithTraits("Elf", "Thief"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
