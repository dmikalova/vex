//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// KompsosHaruspex
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Dinosaur • Priest
//
//	Each friendly creature's play effect is a play/reap effect.
var KompsosHaruspex = card.New(
	"Kompsos Haruspex",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, 224),
	card.WithPower(4),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Priest),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
