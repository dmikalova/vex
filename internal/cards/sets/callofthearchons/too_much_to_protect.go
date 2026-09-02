package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Too Much to Protect
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Steal all but 6 Æmber from your opponent.
var TooMuchToProtect = card.New(
	"Too Much to Protect",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 283),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.StealAember{By: card.AllBut(6)}),
)
