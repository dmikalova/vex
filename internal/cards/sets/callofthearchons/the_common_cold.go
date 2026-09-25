package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// The Common Cold
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Deal 1 damage to each creature. You may destroy each Mars creature.
var TheCommonCold = set.New(
	"The Common Cold",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "336"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{
			Effects: []card.Effect{
				card.DealDamage{
					Amount: 1,
					Target: card.Target.EachCreature,
				},
				card.May{
					Do: card.Destroy{
						Target: card.Target.EachCreature.House(card.Houses.Named(card.House.Mars)),
					},
				},
			},
		}),
)
