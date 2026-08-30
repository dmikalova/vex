//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// BearFlute
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: Fully heal an Ancient Bear. If there are no Ancient Bears in play, search your deck and discard pile and put each Ancient Bear from them into your hand. If you do, shuffle your discard pile into your deck.
var BearFlute = card.New(
	"Bear Flute",
	card.House.Untamed,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 340),
	card.WithTraits("Item"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
