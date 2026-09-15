package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Universal Recycle Bin
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Item
//
//	Action: Archive a purged card you own.
var UniversalRecycleBin = set.New(
	"Universal Recycle Bin",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "179"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.Action, card.ArchivePurgedCard{}),
)
