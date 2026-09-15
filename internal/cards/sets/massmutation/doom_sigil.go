//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// DoomSigil
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Power
//
//	Each creature gains poison.
//	If there are no creatures in play, destroy Doom Sigil.
var DoomSigil = set.New(
	"Doom Sigil",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.MM, "277"),
	card.WithTraits(card.Traits.Power),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
