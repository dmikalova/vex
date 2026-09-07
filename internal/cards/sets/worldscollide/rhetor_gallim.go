//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// RhetorGallim
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Dinosaur • Philosopher
//
//	Play: Your opponent's keys cost +3A during their next turn.
//	Reap: You may exalt Rhetor Gallim. If you do, your opponent's keys cost +3A during their next turn.
var RhetorGallim = card.New(
	"Rhetor Gallim",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, 192),
	card.WithPower(3),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Philosopher),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
