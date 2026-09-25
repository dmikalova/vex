package massmutation

import "github.com/dmikalova/vex/internal/card"

// "Borrow"
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Take control of an enemy artifact. It belongs to house Shadows.
var Borrow = set.New(
	"\"Borrow\"",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "263"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.TakeControl{
				Target:   card.Target.EnemyArtifact,
				Duration: card.Duration.UntilCardLeavesPlay,
			},
			card.BelongToHouse{
				Target:   card.Target.Triggering,
				House:    card.House.Self,
				Duration: card.Duration.UntilThisLeavesPlay,
			},
		}}),
)
