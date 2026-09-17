package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Rapid Evolution
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Choose a creature - for each Æmber in your pool, give the chosen creature a +1 power counter.
var RapidEvolution = set.New(
	"Rapid Evolution",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "373"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.ChooseCreatureThen{
			Target: card.Target.Creature,
			Then: card.AddPowerCounter{
				Target: card.Target.TheChosenCreature,
				Amount: 1,
				Per:    card.AemberInPool{Player: card.Controller},
			},
		}),
)
