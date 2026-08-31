package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Glorious Few
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: For each creature your opponent controls in excess of you, gain 1 Æmber.
var GloriousFew = card.New(
	"Glorious Few",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 218),
	card.WithAbility(card.Trigger.Play, card.GainAember{
		Player: card.Controller,
		Amount: 1,
		Per:    card.OpponentExcessCreatures{},
	}),
)
