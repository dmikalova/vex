package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Lord Golgotha
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Armor:  2
//	Traits: Knight • Spirit
//
//	Before Fight: Deal 3 damage to each neighbor of the creature Lord Golgotha fights.
var LordGolgotha = card.New(
	"Lord Golgotha",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 252),
	card.WithPower(5),
	card.WithArmor(2),
	card.WithTraits("Knight", "Spirit"),
	card.WithAbility(
		card.Trigger.BeforeFight, card.DealDamage{
			Target: card.Target.CreatureFought.NeighborsOf(),
			Amount: 3,
		}),
)
