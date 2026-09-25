package massmutation

import "github.com/dmikalova/vex/internal/card"

// Particle Sweep
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Choose a creature. If it is a Mutant creature, destroy the chosen creature. Otherwise, deal 2 damage to the chosen creature.
var ParticleSweep = set.New(
	"Particle Sweep",
	card.House.StarAlliance,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "327"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.ChooseCreatureThen{
			Target: card.Target.Creature,
			Then: card.Conditional{
				Cond: card.ItIsOfTrait{Trait: card.Traits.Mutant},
				Then: card.Destroy{Target: card.Target.TheChosenCreature},
				Else: card.DealDamage{
					Amount: 2,
					Target: card.Target.TheChosenCreature,
				},
			},
		}),
)
