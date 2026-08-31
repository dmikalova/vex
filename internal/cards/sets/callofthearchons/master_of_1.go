package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Master of 1
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Demon
//
//	Reap: You may destroy a creature with power 1.
var MasterOf1 = card.New(
	"Master of 1",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 89),
	card.WithPower(4),
	card.WithTraits("Demon"),
	card.WithAbility(
		card.Trigger.Reap, card.May{
			Do: card.Destroy{Target: card.Target.Creature.PowerExactly(1)},
		}),
)

// TODO: make this a templated card in deckgen for master of 1 to 5
