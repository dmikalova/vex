package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Malison
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Demon
//
//	Fight: Move an enemy creature anywhere in its controller's battleline -> if it is on a flank, the chosen creature captures 1 Æmber from your opponent.
var Malison = set.New(
	"Malison",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "80"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Demon),
	card.WithAbility(
		card.Trigger.Fight, card.Then{
			First: card.MoveWithinBattleline{Target: card.Target.EnemyCreature},
			Result: card.Conditional{
				Cond: card.OnFlank{OfIt: true},
				Then: card.CaptureAember{
					Amount: 1,
					Target: card.Target.TheChosenCreature,
					Source: card.Opponent,
				},
			},
		}),
)
