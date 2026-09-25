package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Library of Babble
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Common
//	Traits: Location
//
//	Action: Draw a card.
var LibraryOfBabble = set.New(
	"Library of Babble",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Common,
	card.Provenance(card.CotA, "129"),
	card.WithTraits(card.Traits.Location),
	card.WithAbility(
		card.Trigger.Action, card.Draw{Amount: 1}),
)
