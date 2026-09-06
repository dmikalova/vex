package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Golden Aura
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Choose a creature - fully heal it, and for the remainder of the turn, it belongs to house Sanctum and cannot be dealt damage.
var GoldenAura = card.New(
	"Golden Aura",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, 217),
	card.WithAbility(card.Trigger.Play, card.ChooseCreatureThen{
		Target: card.Target.Creature,
		Then: card.Sequence{Effects: []card.Effect{
			card.Heal{Fully: true, Target: card.Target.Triggering},
			card.ForDuration{
				Duration: card.Duration.EndOfTurn,
				Effects: []card.Effect{
					card.BelongToHouse{
						Target:   card.Target.Triggering,
						House:    card.House.Self,
						Duration: card.Duration.EndOfTurn,
					},
					card.PreventDamage{
						Target:   card.Target.Triggering,
						Duration: card.Duration.EndOfTurn,
					},
				},
			},
		}},
	}),
)
