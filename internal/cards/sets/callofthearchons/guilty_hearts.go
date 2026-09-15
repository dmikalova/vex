package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Guilty Hearts
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Destroy each creature with Æmber on it.
var GuiltyHearts = set.New(
	"Guilty Hearts",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "61"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Destroy{
			Target: card.Target.EachCreature.WithAember(),
		}),
)
