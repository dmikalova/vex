package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Sic Semper Tyrannosaurus
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Move all Æmber from the most powerful creature to your pool. Destroy the chosen creature.
var SicSemperTyrannosaurus = set.New(
	"Sic Semper Tyrannosaurus",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "209"),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{
			Effects: []card.Effect{
				card.MoveAember{
					All:  true,
					From: card.Target.EachCreature.Refine(card.MostPowerful),
					To:   card.Controller,
					Bind: true,
				},
				card.Destroy{Target: card.Target.TheChosenCreature},
			},
		}),
)
