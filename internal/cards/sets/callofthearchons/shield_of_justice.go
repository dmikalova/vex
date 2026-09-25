package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Shield of Justice
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: For the remainder of the turn, each friendly creature cannot be dealt damage.
var ShieldOfJustice = set.New(
	"Shield of Justice",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, "225"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.CannotBeDealtDamage{
			Target:   card.Target.EachFriendlyCreature,
			Duration: card.Duration.RemainderOfPlayerTurn,
		}),
)
