package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Pawn Sacrifice
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Destroy a friendly creature -> deal 3 damage to a creature and 3 damage to a different creature.
var PawnSacrifice = set.New(
	"Pawn Sacrifice",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, "279"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Then{
			First: card.Destroy{Target: card.Target.FriendlyCreature},
			Result: card.DealDamage{
				Spread: card.DifferentCreatures{
					First:  3,
					Second: 3,
				},
			},
		}),
)
