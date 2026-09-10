package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Tribute
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: The most powerful friendly creature captures 2 Æmber from your opponent. You may exalt the chosen creature to repeat the preceding effect.
var Tribute = card.New(
	"Tribute",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "196"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.ExaltToRepeat{
			Do: card.CaptureAember{
				Amount: 2,
				Target: card.Target.EachFriendlyCreature.
					Refine(card.MostPowerful(1)),
				Source: card.Opponent,
			},
			Exalt: card.Target.TheChosenCreature,
		}),
)
