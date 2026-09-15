//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// AutoVac5150
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: You may discard a card from your archives. If you do, keys cost +3A during your opponent's next turn. Otherwise, archive a card.
var AutoVac5150 = set.New(
	"Auto-Vac 5150",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.MM, "101"),
	card.WithTraits(card.Traits.Item),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
