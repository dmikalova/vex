package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Guardian Demon
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Demon
//
//	Play/Fight/Reap: Heal 2 damage from a creature. Deal that amount of damage to another creature.
var GuardianDemon = set.New(
	"Guardian Demon",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "88"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Demon),
	card.WithAbility(card.Trigger.PlayFightReap, card.Sequence{Effects: []card.Effect{
		card.Heal{
			Amount: 2,
			Target: card.Target.Creature,
		},
		card.DealDamage{
			AmountFrom: card.DamageHealed{},
			Target:     card.Target.OtherCreature,
		},
	}}),
)
