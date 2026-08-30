package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Guilty Hearts
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Destroy each creature with Æmber on it.
var GuiltyHearts = card.New(
	"Guilty Hearts",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 61),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Destroy{
			Target: card.Target.EachCreature.WithAember(),
		}),
)
