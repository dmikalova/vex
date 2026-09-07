//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Philophosaurus
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Dinosaur • Philosopher
//
//	Reap: You may look at the top 3 cards of your deck. Archive 1, add 1 to your hand, and discard 1.
var Philophosaurus = card.New(
	"Philophosaurus",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "207"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Philosopher),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
