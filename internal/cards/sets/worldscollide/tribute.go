package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Tribute
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: The most powerful friendly creature captures 2 Æmber from your opponent. You may exalt the chosen creature to repeat the preceding effect.
var Tribute = set.New(
	"Tribute",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "196"),
	card.WithBonus(card.Bonus.Aember),
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
