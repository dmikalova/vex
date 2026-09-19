package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Boosted B4-RRY
//
//	House:  Shadows
//	Type:   Gigantic Creature
//	Rarity: Special
//	Power:  7
//	Armor:  2
//	Traits: Robot
//
//	Play/Fight/Reap: Choose one:
//	- Take control of an enemy artifact. If it does not belong to a house on your identity, it belongs to house Shadows.
//	- Play a random card from your opponent's archives.
var BoostedB4RRY = set.Gigantic(
	"Boosted B4-RRY",
	card.House.Shadows,
	card.Rarity.Special,
	card.Provenance(card.MoMu, "257"),
	card.WithPower(7),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Robot),
	card.WithAbility(
		card.Trigger.PlayFightReap, card.ChooseOne{
			Options: []card.Effect{
				card.Sequence{Effects: []card.Effect{
					card.TakeControl{
						Target:   card.Target.EnemyArtifact,
						Duration: card.Duration.UntilCardLeavesPlay,
					},
					card.Conditional{
						Cond: card.ItIsOffIdentity{},
						Then: card.BelongToHouse{
							Target:   card.Target.Triggering,
							House:    card.House.Self,
							Duration: card.Duration.UntilThisLeavesPlay,
						},
					},
				}},
				card.PlayFromOpponent{From: card.Archives},
			},
		}),
)
