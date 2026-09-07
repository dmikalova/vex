package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Irestaff
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Common
//	Æmber:  1
//	Traits: Weapon
//
//	Action: Choose a creature - enrage it, and give it a +1 power counter.
var Irestaff = card.New(
	"Irestaff",
	card.House.Brobnar,
	card.Type.Artifact,
	card.Rarity.Common,
	card.Provenance(card.WC, "10"),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Weapon),
	card.WithAbility(
		card.Trigger.Action, card.ChooseCreatureThen{
			Target: card.Target.Creature,
			Then: card.Sequence{Effects: []card.Effect{
				card.Enrage{Target: card.Target.Triggering},
				card.AddPowerCounter{
					Target: card.Target.Triggering,
					Amount: 1,
				},
			}},
		}),
)
