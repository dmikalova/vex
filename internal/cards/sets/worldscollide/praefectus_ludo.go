//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// PraefectusLudo
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Dinosaur • Politician
//
//	Each other friendly creature gains, "Destroyed: Move each A on this creature to the common supply."
var PraefectusLudo = card.New(
	"Praefectus Ludo",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, 190),
	card.WithPower(5),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Politician),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
