package massmutation

import "github.com/dmikalova/vex/internal/card"

// Burning Glare
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Choose one:
//	- Stun an enemy creature
//	- Stun each enemy Mutant creature.
//	Enhance Damage.
var BurningGlare = set.New(
	"Burning Glare",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "128"),
	card.WithBonus(card.Bonus.Aember),
	card.WithEnhance(card.Bonus.Damage),
	card.WithAbility(
		card.Trigger.Play, card.ChooseOne{
			Options: []card.Effect{
				card.Stun{Target: card.Target.EnemyCreature},
				card.Stun{
					Target: card.Target.EachEnemyCreature.WithTrait(card.Traits.Mutant),
				},
			},
		}),
)
