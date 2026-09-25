package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Interdimensional Graft
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: During your opponent's next turn, after forging a key, your opponent gives you all their Æmber.
var InterdimensionalGraft = set.New(
	"Interdimensional Graft",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "112"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.ForOpponentNextTurn{
			On: card.Event.Forge,
			Do: card.GiveAember{All: true},
		}),
)
