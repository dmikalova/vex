package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// The Howling Pit
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	During their "draw cards" step, each player refills their hand to 1 additional card.
var TheHowlingPit = card.New(
	"The Howling Pit",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 135),
	card.WithTraits(card.Traits.Location),
	card.WithDrawModifier(card.EachPlayer, 1),
)
