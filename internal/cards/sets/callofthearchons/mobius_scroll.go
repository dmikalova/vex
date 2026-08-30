//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// MobiusScroll
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: Archive Mobius Scroll and up to 2 cards from your hand.
var MobiusScroll = card.New(
	"Mobius Scroll",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 130),
	card.WithTraits("Item"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
