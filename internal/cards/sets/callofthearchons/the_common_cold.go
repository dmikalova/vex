package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// The Common Cold
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Deal 1 damage to each creature. You may destroy each Mars creature.
var TheCommonCold = card.New(
	"The Common Cold",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 336),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{
			Effects: []card.Effect{
				card.Sentence{Effect: card.DealDamage{
					Amount: 1,
					Target: card.Target.EachCreature,
				}},
				card.Sentence{Effect: card.May{
					Do: card.Destroy{Target: card.Target.EachCreature.OfHouse(card.House.Mars)},
				}},
			},
		}),
)
