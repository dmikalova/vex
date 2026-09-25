package massmutation

import "github.com/dmikalova/vex/internal/card"

// Beware the Ides
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Deal 23 damage to a creature in the center of its controller's battleline.
var BewareTheIdes = set.New(
	"Beware the Ides",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "184"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount: 23,
			Target: card.Target.Creature.InCenter(),
		}),
)
