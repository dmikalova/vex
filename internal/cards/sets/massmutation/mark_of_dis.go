package massmutation

import "github.com/dmikalova/vex/internal/card"

// Mark of Dis
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Deal 2 damage to a creature. If it is not destroyed, its controller must choose that creature's house as their active house during their next turn.
var MarkOfDis = set.New(
	"Mark of Dis",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "011"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount: 2,
			Target: card.Target.Creature,
			After:  card.IfSurvives,
			Then: card.MustChooseHouse{
				Player:    card.ItsController,
				Reference: card.ItActiveHouse,
			},
		}),
)
