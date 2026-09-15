package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// King of the Crag
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  7
//	Traits: Giant
//
//	Each enemy Brobnar creature gains -2 power.
var KingOfTheCrag = set.New(
	"King of the Crag",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "38"),
	card.WithPower(7),
	card.WithTraits(card.Traits.Giant),
	card.WithConstant(card.ConstantAbility{
		PowerBonus: -2,
		Target:     card.Target.EachEnemyCreature.House(card.Houses.Named(card.House.Self)),
	}),
)
