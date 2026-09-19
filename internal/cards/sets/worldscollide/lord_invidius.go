package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Lord Invidius
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Traits: Demon • Leader
//
//	Elusive.
//	While Lord Invidius is in the center of your battleline, it gains, "Reap: Take control of an enemy flank creature and exhaust it. It belongs to house Dis."
var LordInvidius = set.New(
	"Lord Invidius",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "110"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Demon, card.Traits.Leader),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithConstant(card.ConstantAbility{
		WhileInCenter: true,
		Target:        card.Target.This,
		Granted: []card.Ability{{
			Trigger: card.Trigger.Reap,
			Effect: card.Sequence{Effects: []card.Effect{
				card.TakeControl{
					Target:     card.Target.EnemyCreature.OnFlank(),
					Duration:   card.Duration.UntilCardLeavesPlay,
					AndExhaust: true,
				},
				card.BelongToHouse{
					Target:   card.Target.Triggering,
					House:    card.House.Self,
					Duration: card.Duration.UntilThisLeavesPlay,
				},
			}},
		}},
	}),
)
