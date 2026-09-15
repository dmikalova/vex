package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// GySgt. Margot
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Mutant
//
//	Fight/Reap: Deal 2 damage to an enemy creature. Ward a friendly creature.
var GySgtMargot = set.New(
	"GySgt. Margot",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MoMu, "360"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Mutant),
	card.WithAbility(
		card.Trigger.FightReap, card.Sentences{
			Effects: []card.Effect{
				card.DealDamage{Amount: 2, Target: card.Target.EnemyCreature},
				card.Ward{Target: card.Target.FriendlyCreature},
			},
		}),
)
