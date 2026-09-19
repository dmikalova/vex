package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Quintrino Flux
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Choose a friendly creature and an enemy creature. Destroy each creature with the same power as either of the chosen creatures.
var QuintrinoFlux = set.New(
	"Quintrino Flux",
	card.House.StarAlliance,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "317"),
	card.WithAbility(
		card.Trigger.Play,
		card.Destroy{
			Target: card.Target.EachCreature.Refine(card.SamePowerAsEitherChosen),
		},
	),
)
