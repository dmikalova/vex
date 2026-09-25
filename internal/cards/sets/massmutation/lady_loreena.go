package massmutation

import "github.com/dmikalova/vex/internal/card"

// Lady Loreena
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Armor:  3
//	Traits: Spirit • Knight
//
//	Taunt.
//	Lady Loreena's taunt also applies to its neighbors' neighbors.
var LadyLoreena = set.New(
	"Lady Loreena",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "165"),
	card.WithPower(6),
	card.WithArmor(3),
	card.WithTraits(card.Traits.Spirit, card.Traits.Knight),
	card.WithKeywords(card.Keyword.Taunt),
	card.WithTauntReachingNeighborsNeighbors(),
)
