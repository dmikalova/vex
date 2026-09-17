package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Shadowsaurus
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Armor:  1
//	Traits: Dinosaur • Thief
//
//	Action: Move all Æmber from an enemy creature to your opponent's pool, and if you moved any Æmber this way, take control of it, and that creature belongs to house Shadows.
var Shadowsaurus = set.New(
	"Shadowsaurus",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "292"),
	card.WithPower(5),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Thief),
	card.WithAbility(
		card.Trigger.Action, card.Sequence{Effects: []card.Effect{
			card.MoveAember{
				All:  true,
				From: card.Target.EnemyCreature,
				To:   card.Opponent,
				Bind: true,
			},
			card.Conditional{
				Cond: card.MovedAnyAember{},
				Then: card.Sequence{Effects: []card.Effect{
					card.TakeControl{
						Target:   card.Target.Triggering,
						Duration: card.Duration.Forever,
					},
					card.BelongToHouse{
						Target:   card.Target.Triggering,
						House:    card.House.Shadows,
						Duration: card.Duration.UntilThisLeavesPlay,
						Pronoun:  true,
					},
				}},
			},
		}}),
)
