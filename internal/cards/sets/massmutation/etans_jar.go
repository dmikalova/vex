//go:build todo

package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// EtansJar
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Play: Name a card. Until Etan's Jar leaves play, cards with that name cannot be played.
var EtansJar = set.New(
	"Etan's Jar",
	card.House.Dis,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.MM, "035"),
	card.WithTraits(card.Traits.Item),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
