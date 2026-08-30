//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// NiffleApe
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Beast • Niffle
//
//	While Niffle Ape is attacking, ignore taunt and elusive.
var NiffleApe = card.New(
	"Niffle Ape",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 363),
	card.WithPower(3),
	card.WithTraits("Beast", "Niffle"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
