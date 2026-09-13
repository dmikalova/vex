package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Quintrino Flux
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Choose a friendly Creature and an enemy Creature - destroy each Creature with the same power as either of the chosen Creatures.
var QuintrinoFlux = card.New(
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
