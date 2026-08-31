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
//	Reap: You may fully heal a creature -> for the remainder of the turn, it cannot be dealt damage.
var Protectrix = card.New(
	"Protectrix",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 254),
	card.WithPower(5),
	card.WithTraits("Knight", "Spirit"),
	card.WithAbility(card.Trigger.Reap, card.May{
		Do: card.Then{
			First: card.Heal{Fully: true, Target: card.Target.Creature},
			Result: card.PreventDamage{
				Target:   card.Target.Triggering,
				Duration: card.Duration.EndOfTurn,
			},
		},
	}),
)
