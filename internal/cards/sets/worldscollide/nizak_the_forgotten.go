//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// NizakTheForgotten
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Special
//	Power:  6
//	Traits: Dragon • Psion
//
//	While fighting, Nizak, The Forgotten gains invulnerable. (It cannot be destroyed or dealt damage.)
//	After an enemy creature is destroyed fighting Nizak, The Forgotten, return that creature to its owner's hand.
var NizakTheForgotten = card.New(
	"Nizak, The Forgotten",
	card.House.Brobnar,
	card.Type.Creature,
	// TODO(variant): rarity relabelled from FIXED to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.WC, "A05"),
	card.WithPower(6),
	card.WithTraits(card.Traits.Dragon, card.Traits.Psion),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
