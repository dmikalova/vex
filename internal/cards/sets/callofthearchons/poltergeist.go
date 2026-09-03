package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Poltergeist
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Use an artifact. Destroy it.
var Poltergeist = card.New(
	"Poltergeist",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 69),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{Effects: []card.Effect{
			card.Use{
				Max:    1,
				Target: card.Target.EachArtifact,
			},
			card.Destroy{Target: card.Target.Triggering},
		}}),
)
