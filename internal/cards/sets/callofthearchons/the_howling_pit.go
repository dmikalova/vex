//go:build todo

package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// TheHowlingPit
//
// TODO(stub): unimplemented. Remove the //go:build todo tag and
// implement the ability once the needed effect exists.
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Location
//
//	During their "draw cards" step, each player refills their hand to 1 additional card.
var TheHowlingPit = card.New(
	"The Howling Pit",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 135),
	card.WithAemberBonus(1),
	card.WithTraits("Location"),
	// TODO(stub): add WithKeywords / WithAbility for the printed text above.
)
