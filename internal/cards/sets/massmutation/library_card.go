//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// LibraryCard
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: Purge Library Card. If you do, for the remainder of the turn, after you play a card, draw a card.
var LibraryCard = set.New(
	"Library Card",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.MM, "105"),
	card.WithTraits(card.Traits.Item),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
