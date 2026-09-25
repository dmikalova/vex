package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Memory Chip
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Item
//
//	After you choose Logos as your active house, archive a card from your hand.
var MemoryChip = set.New(
	"Memory Chip",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "132"),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.AfterChooseHouse, card.Conditional{
			Cond: card.ChoseHouse{House: card.House.Self},
			Then: card.ArchiveCard{
				Zone:      card.Hand,
				Selection: card.Chosen{},
			},
		}),
)
