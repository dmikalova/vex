//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// AutoEncoder
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Common
//	Traits: Item
//
//	After a card is discarded from your hand, archive the top card of your deck.
var AutoEncoder = set.New(
	"Auto-Encoder",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Common,
	card.Provenance(card.MM, "066"),
	card.WithTraits(card.Traits.Item),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
