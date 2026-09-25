package massmutation

import "github.com/dmikalova/vex/internal/card"

// Hadron Collision
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Choose a creature. Remove a ward from the chosen creature. Deal 3 damage to the chosen creature, ignoring armor.
var HadronCollision = set.New(
	"Hadron Collision",
	card.House.StarAlliance,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "308"),
	card.WithAbility(
		card.Trigger.Play, card.ChooseCreatureThen{
			Target: card.Target.Creature,
			Then: card.Sequence{Effects: []card.Effect{
				card.RemoveWard{Target: card.Target.TheChosenCreature},
				card.DealDamage{
					Amount:      3,
					Target:      card.Target.TheChosenCreature,
					IgnoreArmor: true,
				},
			}},
		}),
)
