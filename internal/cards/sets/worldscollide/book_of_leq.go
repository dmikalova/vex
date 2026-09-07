//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// BookOfLeQ
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: Reveal the top card of your deck. If it is a non-Star Alliance card, its house becomes your active house. Otherwise, end your turn.
var BookOfLeQ = card.New(
	"Book of leQ",
	card.House.Staralliance,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, 325),
	card.WithTraits(card.Traits.Item),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
