package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Miasma
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Your opponent skips the "forge a key" phase during their next turn.
var Miasma = set.New(
	"Miasma",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, "275"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(card.Trigger.Play, card.SkipForgePhase{Player: card.Opponent}),
)
