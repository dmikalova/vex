package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Pawn Sacrifice
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Destroy a friendly Creature -> deal 3 damage to a Creature and deal 3 damage to a different Creature.
var PawnSacrifice = set.New(
	"Pawn Sacrifice",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, "279"),
	card.WithAemberBonus(1),
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
