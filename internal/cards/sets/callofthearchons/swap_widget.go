//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// SwapWidget
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: Return a ready friendly Mars creature to your hand. If you do, put a Mars creature with a different name from your hand into play, then ready it.
var SwapWidget = card.New(
	"Swap Widget",
	card.House.Mars,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 189),
	card.WithTraits("Item"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
