package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Pandemonium
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Each undamaged creature captures 1 Æmber from its opponent.
var Pandemonium = card.New(
	"Pandemonium",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 68),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.CaptureAember{
			Amount: 1,
			Target: card.Target.EachCreature.Undamaged(),
			Source: card.ItsOpponent,
		}),
)
