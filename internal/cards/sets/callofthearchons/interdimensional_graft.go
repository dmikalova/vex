package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Interdimensional Graft
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: During your opponent's next turn, after forging a key, your opponent gives you all their Æmber.
var InterdimensionalGraft = card.New(
	"Interdimensional Graft",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "112"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.ForOpponentNextTurn{
			On: card.Event.Forge,
			Do: card.GiveAember{All: true},
		}),
)
