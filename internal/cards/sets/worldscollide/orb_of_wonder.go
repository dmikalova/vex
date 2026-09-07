//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// OrbOfWonder
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Special
//	Traits: Item
//
//	Omni: Sacrifice Orb of Wonder. If you do, search your deck for a card and add it to your hand. Then, shuffle your deck.
var OrbOfWonder = card.New(
	"Orb of Wonder",
	card.House.Brobnar,
	card.Type.Artifact,
	// TODO(variant): rarity relabelled from FIXED to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.WC, "A06"),
	card.WithTraits(card.Traits.Item),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
