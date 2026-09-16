package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Matter Maker
//
//	House:  Star Alliance
//	Type:   Artifact
//	Rarity: Rare
//	Bonus:  Æmber
//	Traits: Item
//
//	You may play upgrades as if they were in the active house.
var MatterMaker = set.New(
	"Matter Maker",
	card.House.StarAlliance,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.MM, "349"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Item),
	card.WithPlayPermission(card.PlayPermission{
		Types: card.Types.Of(card.Type.Upgrade),
	}),
)
