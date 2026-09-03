package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Red-Hot Armor
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Each enemy creature with armor loses all of its armor. Deal 1 damage to each enemy creature with armor for each point of armor it lost this way.
var RedHotArmor = card.New(
	"Red-Hot Armor",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 70),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{Effects: []card.Effect{
			card.LoseArmor{Target: card.Target.EachEnemyCreature.WithArmor()},
			card.DealDamage{
				Target:    card.Target.EachEnemyCreature.WithArmor(),
				Amount:    1,
				PerTarget: card.ArmorLostThisWay,
			},
		}}),
)
