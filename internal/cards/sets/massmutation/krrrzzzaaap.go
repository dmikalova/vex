package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Krrrzzzaaap!!!
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Destroy each non-Mutant creature. Gain 1 chain.
var Krrrzzzaaap = set.New(
	"Krrrzzzaaap!!!",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "090"),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.Destroy{Target: card.Target.EachCreature.ExceptTrait(card.Traits.Mutant)},
			card.GainChains{Player: card.Controller, Amount: 1},
		}}),
)
