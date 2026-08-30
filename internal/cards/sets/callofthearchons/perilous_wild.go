package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Perilous Wild
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Destroy each elusive creature.
var PerilousWild = card.New(
	"Perilous Wild",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 331),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Destroy{
			Target: card.Target.EachCreature.Keyword(card.Keyword.Elusive),
		}),
)
