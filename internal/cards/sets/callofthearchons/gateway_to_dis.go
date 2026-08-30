package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Gateway to Dis
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Destroy each creature. Gain 3 chains.
var GatewayToDis = card.New(
	"Gateway to Dis",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 59),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{
			Effects: []card.Effect{
				card.Sentence{Effect: card.Destroy{Target: card.Target.EachCreature}},
				card.Sentence{Effect: card.GainChains{Amount: 3}},
			},
		}),
)
