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
//	Rarity: FIXED
//	Traits: Item
//
//	Omni: Sacrifice Orb of Wonder. If you do, search your deck for a card and add it to your hand. Then, shuffle your deck.
var OrbOfWonder = card.New(
	"Orb of Wonder",
	card.House.Brobnar,
	card.Type.Artifact,
	card.Rarity.FIXED,
	card.Provenance(card.WC, 0),
	card.WithTraits(card.Traits.Item),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
