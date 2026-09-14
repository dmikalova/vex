package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Universal Recycle Bin
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Æmber:  1
//	Traits: Item
//
//	Action: Archive a purged card you own.
var UniversalRecycleBin = set.New(
	"Universal Recycle Bin",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "179"),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.Action, card.ArchivePurgedCard{}),
)
