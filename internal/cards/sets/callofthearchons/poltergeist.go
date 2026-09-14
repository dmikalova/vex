package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Poltergeist
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Use an Artifact. Destroy it.
var Poltergeist = set.New(
	"Poltergeist",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "69"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{Effects: []card.Effect{
			card.Use{
				Max:          1,
				Target:       card.Target.EachArtifact,
				EvenUnusable: true,
			},
			card.Destroy{Target: card.Target.Triggering},
		}}),
)
