//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// MimicGel
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Rare
//	Power:  0
//	Traits: Shapeshifter • Mutant
//
//	Mimic Gel cannot be played unless there is another creature in play.
//	Mimic Gel enters play as a copy of another creature in play, except it belongs to house Logos.
var MimicGel = card.New(
	"Mimic Gel",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, 170),
	card.WithPower(0),
	card.WithTraits(card.Traits.Shapeshifter, card.Traits.Mutant),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
