//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// IronObelisk
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Your opponent's keys cost +1 Aember for each friendly damaged Brobnar creature.
var IronObelisk = card.New(
	"Iron Obelisk",
	card.House.Brobnar,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 23),
	card.WithTraits("Location"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
