package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Brain Stem Antenna
//
//	House:  Mars
//	Type:   Upgrade
//	Rarity: Rare
//
//	This creature gains, "After you play a Mars creature, ready this creature, and for the remainder of the turn, this creature belongs to house Mars."
var BrainStemAntenna = card.New(
	"Brain Stem Antenna",
	card.House.Mars,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 209),
	card.WithStatic(card.StaticModifier{
		Granted: []card.Ability{
			{Trigger: card.Trigger.AfterCardPlayed, Effect: card.Conditional{
				Cond: card.ItIs{House: card.House.Self, Type: card.Type.Creature},
				Then: card.Sequence{
					Effects: []card.Effect{
						card.Ready{Target: card.Target.This},
						card.BelongToHouse{
							Target:   card.Target.This,
							House:    card.House.Self,
							Duration: card.Duration.EndOfTurn,
						},
					},
				},
			}},
		},
	}),
)
