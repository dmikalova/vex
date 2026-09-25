package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Pandemonium
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Each undamaged creature captures 1 Æmber from its opponent.
var Pandemonium = set.New(
	"Pandemonium",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "68"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.CaptureAember{
			Amount: 1,
			Target: card.Target.EachCreature.Undamaged(),
			Source: card.ItsOpponent,
		}),
)
