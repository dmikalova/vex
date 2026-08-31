package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Banner of Battle
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Each friendly creature gains +1 power.
var BannerOfBattle = card.New(
	"Banner of Battle",
	card.House.Brobnar,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 20),
	card.WithTraits("Item"),
	card.WithConstant(card.ConstantAbility{
		PowerBonus: 1,
		Target:     card.Target.EachFriendlyCreature,
	}),
)
