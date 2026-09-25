package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Library of the Damned
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Location
//
//	Action: Archive a card from your hand.
var LibraryOfTheDamned = set.New(
	"Library of the Damned",
	card.House.Dis,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "76"),
	card.WithTraits(card.Traits.Location),
	card.WithAbility(
		card.Trigger.Action,
		card.ArchiveCard{
			Zone:      card.Hand,
			Selection: card.Chosen{},
		},
	),
)
