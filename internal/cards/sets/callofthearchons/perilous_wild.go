package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Perilous Wild
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Destroy each elusive creature.
var PerilousWild = set.New(
	"Perilous Wild",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "331"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Destroy{
			Target: card.Target.EachCreature.Keyword(card.Keyword.Elusive),
		}),
)
