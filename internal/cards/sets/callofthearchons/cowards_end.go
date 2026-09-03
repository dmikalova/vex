package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Coward's End
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Destroy each undamaged creature. Gain 3 chains.
var CowardsEnd = card.New(
	"Coward's End",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 7),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{
			Effects: []card.Effect{
				card.Destroy{Target: card.Target.EachCreature.Undamaged()},
				card.GainChains{Amount: 3},
			},
		}),
)
