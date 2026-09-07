//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// OratorHissaro
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Dinosaur • Politician
//
//	Deploy.
//	Play: Ready and exalt each of Orator Hissaro's neighbors. For the remainder of the turn, they belong to house Saurian.
var OratorHissaro = card.New(
	"Orator Hissaro",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "205"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Politician),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
