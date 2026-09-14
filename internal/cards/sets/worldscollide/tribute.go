package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Tribute
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: The most powerful friendly Creature captures 2 Æmber from your opponent. You may exalt the chosen Creature to repeat the preceding effect.
var Tribute = set.New(
	"Tribute",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "196"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Repeat{
			Do: card.CaptureAember{
				Amount: 2,
				Target: card.Target.EachFriendlyCreature.
					Refine(card.MostPowerful),
				Source: card.Opponent,
			},
			Gate: card.ByExalting{Creature: card.Target.TheChosenCreature},
		}),
)
