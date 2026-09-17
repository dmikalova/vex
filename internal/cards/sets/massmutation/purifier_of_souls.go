package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Purifier of Souls
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Armor:  2
//	Traits: Human • Priest
//
//	Destroyed effects cannot trigger.
var PurifierOfSouls = set.New(
	"Purifier of Souls",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "174"),
	card.WithPower(5),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Human, card.Traits.Priest),
	card.WithConstant(card.ConstantAbility{
		DisableTriggers: card.Triggers(card.Trigger.Destroyed),
	}),
)
