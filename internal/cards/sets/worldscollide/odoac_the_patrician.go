//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// OdoacThePatrician
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
//	Play: Capture 1A.
//	While Odoac the Patrician has A on it, your A cannot be stolen.
var OdoacThePatrician = card.New(
	"Odoac the Patrician",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, 188),
	card.WithPower(5),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Politician),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
