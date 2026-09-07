//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// TheColosseum
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	After an enemy creature is destroyed while fighting, put a glory counter on The Colosseum.
//	Omni: If there are 6 or more glory counters on The Colosseum, remove 6 and forge a key at current cost.
var TheColosseum = card.New(
	"The Colosseum",
	card.House.Saurian,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "233"),
	card.WithTraits(card.Traits.Location),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
