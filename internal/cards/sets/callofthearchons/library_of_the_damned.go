package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Library of the Damned
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Location
//
//	Action: Archive a card from your hand.
var LibraryOfTheDamned = card.New(
	"Library of the Damned",
	card.House.Dis,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 76),
	card.WithTraits(card.Traits.Location),
	card.WithAbility(card.Trigger.Action, card.ArchiveFromHand{Count: 1}),
)
