package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Brain Stem Antenna
//
//	House:  Mars
//	Type:   Upgrade
//	Rarity: Rare
//
//	This Creature gains, "After you play a Mars Creature, ready this Creature, and for the remainder of the turn, this Creature belongs to house Mars."
var BrainStemAntenna = set.New(
	"Brain Stem Antenna",
	card.House.Mars,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "209"),
	card.WithStatic(card.StaticModifier{
		Granted: []card.Ability{
			{Trigger: card.Trigger.AfterCardPlayed, Effect: card.Conditional{
				Cond: card.ItIs{
					House: card.Houses.Named(card.House.Self),
					Type:  card.Type.Creature,
				},
				Then: card.Sequence{
					Effects: []card.Effect{
						card.Ready{Target: card.Target.This},
						card.BelongToHouse{
							Target:   card.Target.This,
							House:    card.House.Self,
							Duration: card.Duration.RemainderOfPlayerTurn,
						},
					},
				},
			}},
		},
	}),
)
