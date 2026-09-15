package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Effervescent Principle
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Each player loses half of their Æmber, rounded down, and gain 1 chain.
var EffervescentPrinciple = set.New(
	"Effervescent Principle",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, "109"),
	card.WithoutEnhancement(card.Bonus.Capture),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{
			Effects: []card.Effect{
				card.LoseAember{
					Player: card.EachPlayer,
					By:     card.HalfRoundedDown,
				},
				card.GainChains{Amount: 1},
			},
		}),
)
