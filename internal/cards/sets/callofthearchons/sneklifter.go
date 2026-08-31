package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Sneklifter
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Elf • Thief
//
//	Play: Take control of an enemy artifact. If it does not belong to a house on your identity, it belongs to house Shadows until it leaves play.
var Sneklifter = card.New(
	"Sneklifter",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 313),
	card.WithPower(2),
	card.WithTraits("Elf", "Thief"),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.Sentence{Effect: card.TakeControl{
				Target:   card.Target.EnemyArtifact,
				Duration: card.Duration.Forever,
			}},
			card.Conditional{
				Cond: card.ItIsOffIdentity{},
				Then: card.BelongToHouse{
					Target:   card.Target.Triggering,
					House:    card.House.Shadows,
					Duration: card.Duration.UntilThisLeavesPlay,
				},
			},
		}}),
)
