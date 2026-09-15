package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Dark Minion
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Mutant
//
//	Destroyed: Deal 1 damage to each enemy creature.
//	Enhance Damage.
var DarkMinion = set.New(
	"Dark Minion",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "004"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Mutant),
	card.WithEnhance(card.Bonus.Damage),
	card.WithAbility(
		card.Trigger.Destroyed, card.DealDamage{
			Amount: 1,
			Target: card.Target.EachEnemyCreature,
		}),
)
