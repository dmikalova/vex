package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Irestaff
//
//	House:  Brobnar
//	Type:   Artifact
//	Rarity: Common
//	Bonus:  Æmber
//	Traits: Weapon
//
//	Action: Choose a creature. Enrage it. Give it a +1 power counter.
var Irestaff = set.New(
	"Irestaff",
	card.House.Brobnar,
	card.Type.Artifact,
	card.Rarity.Common,
	card.Provenance(card.WC, "10"),
	card.WithBonus(card.Bonus.Aember),
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
