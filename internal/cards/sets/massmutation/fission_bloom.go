//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// FissionBloom
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Power
//
//	Enhance Draw.
//	Action: The next time you play a card this turn, resolve each of its bonus icons an additional time.
var FissionBloom = set.New(
	"Fission Bloom",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "087"),
	card.WithTraits(card.Traits.Power),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
