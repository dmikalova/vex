//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// HarlandMindlock
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Cyborg • Scientist
//
//	Play: Take control of an enemy flank creature until Harland Mindlock leaves play.
var HarlandMindlock = card.New(
	"Harland Mindlock",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 143),
	card.WithPower(1),
	card.WithTraits("Cyborg", "Scientist"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
