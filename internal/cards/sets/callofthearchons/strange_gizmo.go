//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// StrangeGizmo
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Item
//
//	After you forge a key, destroy each creature and artifact.
var StrangeGizmo = card.New(
	"Strange Gizmo",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 134),
	card.WithAemberBonus(1),
	card.WithTraits("Item"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
