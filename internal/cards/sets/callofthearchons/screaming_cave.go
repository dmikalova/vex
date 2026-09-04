package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Screaming Cave
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Location
//
//	Action: Shuffle your hand and discard pile into your deck.
var ScreamingCave = card.New(
	"Screaming Cave",
	card.House.Dis,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 79),
	card.WithTraits(card.Traits.Location),
	card.WithAbility(
		card.Trigger.Action,
		card.ShuffleIntoDeck{Zones: []card.Zone{card.Hand, card.Discard}},
	),
)
