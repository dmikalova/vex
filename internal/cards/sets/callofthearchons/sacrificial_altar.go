//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// SacrificialAltar
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Location
//
//	Action: Purge a friendly Human creature from play. If you do, play a creature from your discard pile.
var SacrificialAltar = card.New(
	"Sacrificial Altar",
	card.House.Dis,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 78),
	card.WithAemberBonus(1),
	card.WithTraits("Location"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
