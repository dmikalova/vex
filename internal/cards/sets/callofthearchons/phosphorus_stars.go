package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Phosphorus Stars
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Stun each non-Mars Creature. Gain 2 chains.
var PhosphorusStars = set.New(
	"Phosphorus Stars",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, "173"),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{
			Effects: []card.Effect{
				card.Stun{
					Target: card.Target.EachCreature.House(card.Houses.Except(card.House.Self)),
				},
				card.GainChains{Amount: 2},
			},
		}),
)
