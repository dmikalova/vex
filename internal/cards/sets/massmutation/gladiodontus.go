package massmutation

import "github.com/dmikalova/vex/internal/card"

// Gladiodontus
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  15
//	Traits: Mutant
//
//	Gladiodontus deals 5 damage when fighting.
//	Gladiodontus enters play stunned.
//	Fight/Reap: If this is the first time Gladiodontus has been used this turn, ready and enrage Gladiodontus.
var Gladiodontus = set.New(
	"Gladiodontus",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "206"),
	card.WithPower(15),
	card.WithTraits(card.Traits.Mutant),
	card.WithEntersPlay(card.Stun{Target: card.Target.This}),
	card.WithAttackDamage(card.AttackDamage{
		Amount: 5,
		Fixed:  true,
	}),
	card.WithAbility(card.Trigger.FightReap, card.Conditional{
		Cond: card.SourceFirstUseThisTurn{},
		Then: card.Sequence{Effects: []card.Effect{
			card.Ready{Target: card.Target.This},
			card.Enrage{Target: card.Target.This},
		}},
	}),
)
