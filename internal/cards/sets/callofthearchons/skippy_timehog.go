//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// SkippyTimehog
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Mutant
//
//	Play: Your opponent cannot use any cards next turn. (Cards can still be played and discarded.)
var SkippyTimehog = card.New(
	"Skippy Timehog",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 152),
	card.WithPower(1),
	card.WithTraits("Mutant"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
