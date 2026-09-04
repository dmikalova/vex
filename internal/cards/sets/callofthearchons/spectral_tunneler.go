//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// SpectralTunneler
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Uncommon
//	Æmber:  1
//	Traits: Item
//
//	Action: Choose a creature. For the remainder of the turn, that creature is considered a flank creature and gains, "Reap: Draw a card."
var SpectralTunneler = card.New(
	"Spectral Tunneler",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 133),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Item),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
