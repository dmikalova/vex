//go:build todo

package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// FilaTheResearcher
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  1
//	Traits: Human • Scientist
//
//	Elusive. (The first time this creature is attacked each turn, no damage is dealt.)
//	After a creature is played adjacent to Fila the Researcher, draw a card.
var FilaTheResearcher = card.New(
	"Fila the Researcher",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 129),
	card.WithPower(1),
	card.WithTraits(card.Traits.Human, card.Traits.Scientist),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
