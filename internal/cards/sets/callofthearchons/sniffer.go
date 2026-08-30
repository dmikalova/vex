//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Sniffer
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Ally
//
//	Action: For the remainder of the turn, each creature loses elusive.
var Sniffer = card.New(
	"Sniffer",
	card.House.Mars,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 188),
	card.WithAemberBonus(1),
	card.WithTraits("Ally"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
