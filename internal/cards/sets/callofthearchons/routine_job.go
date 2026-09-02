package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Routine Job
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Steal 1 Æmber, and for each copy of Routine Job in your discard pile, steal 1 Æmber.
var RoutineJob = card.New("Routine Job",
	card.House.Shadows, card.Type.Tactic, card.Rarity.Rare,
	card.Provenance(card.CotA, 282),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.StealAember{Amount: 1},
			card.StealAember{
				Amount: 1,
				Per:    card.CopiesInDiscard{},
			},
		}}),
)
