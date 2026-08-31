package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Grey Monk
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Human • Priest
//
//	Each friendly creature gains +1 armor.
//	Reap: Heal 2 damage from a creature.
var GreyMonk = card.New(
	"Grey Monk",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 244),
	card.WithPower(3),
	card.WithTraits("Human", "Priest"),
	card.WithConstant(card.ConstantAbility{
		ArmorBonus: 1,
		Target:     card.Target.EachFriendlyCreature,
	}),
	card.WithAbility(
		card.Trigger.Reap, card.Heal{
			Amount: 2,
			Target: card.Target.Creature,
		}),
)
