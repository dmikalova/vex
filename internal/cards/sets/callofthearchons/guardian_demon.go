package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Guardian Demon
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Demon
//
//	Play/Fight/Reap: Heal 2 damage from a creature. Deal that amount of damage to another creature.
var GuardianDemon = card.New(
	"Guardian Demon",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 88),
	card.WithPower(4),
	card.WithTraits("Demon"),
	card.WithPlayFightReap(card.Sentences{Effects: []card.Effect{
		card.Heal{Amount: 2, Target: card.Target.Creature},
		card.DealDamage{AmountFrom: card.DamageHealed{}, Target: card.Target.OtherCreature},
	}}),
)
