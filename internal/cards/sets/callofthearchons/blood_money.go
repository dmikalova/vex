package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Blood Money
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Exalt an enemy creature 2 times.
var BloodMoney = card.New(
	"Blood Money",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 3),
	card.WithAbility(
		card.Trigger.Play, card.Exalt{
			Target: card.Target.EnemyCreature,
			Amount: 2,
		}),
)
