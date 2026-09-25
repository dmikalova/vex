package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Universal Recycle Bin
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Item
//
//	Action: Archive a card from your purge pile.
var UniversalRecycleBin = set.New(
	"Universal Recycle Bin",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "179"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.Action, card.ArchiveCard{
			Zone:      card.Purged,
			Selection: card.Chosen{},
		}),
)
