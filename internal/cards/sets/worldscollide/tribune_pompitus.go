//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// TribunePompitus
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Armor:  2
//	Traits: Dinosaur • Politician
//
//	Each friendly creature gets +2 power for each A on it.
//	Before Fight: You may exalt Tribune Pompitus.
var TribunePompitus = card.New(
	"Tribune Pompitus",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 213),
	card.WithPower(4),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Politician),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
