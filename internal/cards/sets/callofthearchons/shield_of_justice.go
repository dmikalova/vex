package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Shield of Justice
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: For the remainder of the turn, each friendly creature cannot be dealt damage.
var ShieldOfJustice = card.New(
	"Shield of Justice",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 225),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.PreventDamage{
			Target:   card.Target.EachFriendlyCreature,
			Duration: card.Duration.EndOfTurn,
		}),
)
