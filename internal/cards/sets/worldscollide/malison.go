package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Malison
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Demon
//
//	Fight: You may move an enemy creature anywhere in its controller's battleline -> if it is on a flank, the chosen creature captures 1 Æmber from your opponent.
var Malison = card.New(
	"Malison",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "80"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Demon),
	card.WithAbility(
		card.Trigger.Fight, card.May{Do: card.Then{
			First: card.MoveWithinBattleline{Target: card.Target.EnemyCreature},
			Result: card.Conditional{
				Cond: card.ItIsOnFlank{},
				Then: card.CaptureAember{
					Amount: 1,
					Target: card.Target.TheChosenCreature,
					Source: card.Opponent,
				},
			},
		}},
	),
)
