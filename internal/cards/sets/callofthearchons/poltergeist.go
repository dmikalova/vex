package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Poltergeist
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Use an artifact. Destroy it.
var Poltergeist = set.New(
	"Poltergeist",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "69"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.Use{
				Max:          1,
				Target:       card.Target.EachArtifact,
				EvenUnusable: true,
			},
			card.Destroy{Target: card.Target.Triggering},
		}}),
)
