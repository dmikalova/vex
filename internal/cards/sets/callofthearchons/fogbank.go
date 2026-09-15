package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Fogbank
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Your opponent cannot use creatures to fight during their next turn.
var Fogbank = set.New(
	"Fogbank",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "322"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Restrict{
			Player:   card.Opponent,
			Action:   card.Restricted.Fighting,
			Duration: card.Duration.OpponentNextTurn,
		}),
)
