//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// UnchartedLands
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Staralliance
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Play: Place 6A from the common supply on Uncharted Lands.
//	Each Star Alliance creature gains, "Reap: Move 1A from Uncharted Lands to your pool."
var UnchartedLands = card.New(
	"Uncharted Lands",
	card.House.Staralliance,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, 342),
	card.WithTraits(card.Traits.Location),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
