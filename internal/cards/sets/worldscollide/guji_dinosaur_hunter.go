//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// GujiDinosaurHunter
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Giant • Hunter
//
//	Elusive.
//	Action: Deal 2D to a creature. Deal 6D instead if it is a Dinosaur creature or has A on it.
var GujiDinosaurHunter = card.New(
	"Guji Dinosaur Hunter",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "38"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Giant, card.Traits.Hunter),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
