package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Axiom of Grisk
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Ward a creature. Destroy each creature with no Æmber on it. Gain 2 chains.
var AxiomOfGrisk = set.New(
	"Axiom of Grisk",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "182"),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.Ward{Target: card.Target.Creature},
			card.Destroy{Target: card.Target.EachCreature.WithoutAember()},
			card.GainChains{Amount: 2},
		}}),
)
