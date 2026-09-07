//go:build todo

package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// LibraryOfPolliasaurus
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Saurian
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Location
//
//	Action: Move 1A from a friendly creature to your pool.
var LibraryOfPolliasaurus = card.New(
	"Library of Polliasaurus",
	card.House.Saurian,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 204),
	card.WithTraits(card.Traits.Location),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
