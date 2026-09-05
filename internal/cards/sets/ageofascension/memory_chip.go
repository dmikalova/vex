package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Memory Chip
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Item
//
//	Archive a card from your hand.
var MemoryChip = card.New(
	"Memory Chip",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 132),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.AfterChooseHouse, card.ArchiveFromHand{Count: 1}),
)
