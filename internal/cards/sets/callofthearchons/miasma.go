package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Miasma
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Your opponent skips the "forge a key" step during their next turn.
var Miasma = card.New(
	"Miasma",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 275),
	card.WithAemberBonus(1),
	card.WithAbility(card.Trigger.Play, card.SkipForgeStep{Player: card.Opponent}),
)
