//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Keyforgery
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Item
//
//	When your opponent would forge a key, that player names a house. Reveal a random card from your hand. If that card is not of the named house, destroy Keyforgery and they do not forge that key (no A is spent).
var Keyforgery = card.New(
	"Keyforgery",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, 271),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Item),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
