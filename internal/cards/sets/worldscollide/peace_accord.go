//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// PeaceAccord
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Law
//
//	Play: Each player gains 2A.
//	After a creature is used to fight, its controller loses 4A. Destroy Peace Accord.
var PeaceAccord = card.New(
	"Peace Accord",
	card.House.Staralliance,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "335"),
	card.WithTraits(card.Traits.Law),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
