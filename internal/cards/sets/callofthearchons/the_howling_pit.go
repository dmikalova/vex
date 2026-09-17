package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// The Howling Pit
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Each player's hand size is 1 more.
var TheHowlingPit = set.New(
	"The Howling Pit",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "135"),
	card.WithTraits(card.Traits.Location),
	card.WithDrawModifier(card.EachPlayer, 1),
)
