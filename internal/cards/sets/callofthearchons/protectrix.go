package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Protectrix
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Knight • Spirit
//
//	Reap: Choose a creature - fully heal it, and for the remainder of the turn, it cannot be dealt damage.
var Protectrix = card.New(
	"Protectrix",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 254),
	card.WithPower(5),
	card.WithTraits(card.Traits.Knight, card.Traits.Spirit),
	card.WithAbility(card.Trigger.Reap, card.ChooseCreatureThen{
		Target: card.Target.Creature,
		Then: card.Sequence{Effects: []card.Effect{
			card.Heal{Fully: true, Target: card.Target.Triggering},
			card.PreventDamage{
				Target:   card.Target.Triggering,
				Duration: card.Duration.EndOfTurn,
			},
		}},
	}),
)
