package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Fogbank
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Your opponent cannot use creatures to fight during their next turn.
var Fogbank = card.New(
	"Fogbank",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 110),
	card.Provenance(card.CotA, 322),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.CannotFight{
			Player:   card.Opponent,
			Duration: card.Duration.NextTurn,
		}),
)
