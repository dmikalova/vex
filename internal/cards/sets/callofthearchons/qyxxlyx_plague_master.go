package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Qyxxlyx Plague Master
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Martian • Scientist
//
//	Fight/Reap: Deal 3 damage to each Human creature, ignoring armor.
var QyxxlyxPlagueMaster = set.New(
	"Qyxxlyx Plague Master",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "198"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Martian, card.Traits.Scientist),
	card.WithAbility(card.Trigger.FightReap, card.DealDamage{
		Amount:      3,
		Target:      card.Target.EachCreature.WithTrait(card.Traits.Human),
		IgnoreArmor: true,
	}),
)
