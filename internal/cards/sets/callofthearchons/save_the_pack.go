package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Save the Pack
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Destroy each damaged creature. Gain 1 chain.
var SaveThePack = card.New(
	"Save the Pack",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 333),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{
			Effects: []card.Effect{
				card.Destroy{Target: card.Target.EachCreature.Damaged()},
				card.GainChains{Amount: 1},
			},
		}),
)
