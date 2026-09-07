//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// SciOfficerMorpheus
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Armor:  1
//	Traits: Shapeshifter • Scientist
//
//	After you play a creature with a play effect, trigger its play effect an additional time.
var SciOfficerMorpheus = card.New(
	"Sci. Officer Morpheus",
	card.House.Staralliance,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 318),
	card.WithPower(2),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Shapeshifter, card.Traits.Scientist),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
