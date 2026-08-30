//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// IncubationChamber
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Mars
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Omni: Reveal a Mars creature from your hand. If you do, archive it.
var IncubationChamber = card.New(
	"Incubation Chamber",
	card.House.Mars,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 186),
	card.WithTraits("Location"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
