package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Phosphorus Stars
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Stun each non-Mars creature. Gain 2 chains.
var PhosphorusStars = card.New(
	"Phosphorus Stars",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 173),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{
			Effects: []card.Effect{
				card.Sentence{
					Effect: card.Stun{
						Target: card.Target.EachCreature.ExceptHouse(card.House.Mars),
					},
				},
				card.Sentence{Effect: card.GainChains{Amount: 2}},
			},
		}),
)
