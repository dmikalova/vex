package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Dark Wave
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Deal 2 damage to each non-Mutant creature.
var DarkWave = set.New(
	"Dark Wave",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "247"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount: 2,
			Target: card.Target.EachCreature.ExceptTrait(card.Traits.Mutant),
		}),
)
