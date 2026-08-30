//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// SoulSnatcher
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Vehicle
//
//	Each time a creature is destroyed, its owner gains 1 Aember.
var SoulSnatcher = card.New(
	"Soul Snatcher",
	card.House.Dis,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 80),
	card.WithTraits("Vehicle"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
