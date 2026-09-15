package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Too Much to Protect
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Steal all but 6 Æmber from your opponent.
var TooMuchToProtect = set.New(
	"Too Much to Protect",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "283"),
	card.OneCopyPerDeck(),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.StealAember{By: card.AllBut(6)}),
)
