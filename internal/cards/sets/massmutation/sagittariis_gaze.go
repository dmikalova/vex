package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Sagittarii's Gaze
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Exalt a damaged creature.
//	Enhance Damage.
var SagittariisGaze = set.New(
	"Sagittarii's Gaze",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "196"),
	card.WithBonus(card.Bonus.Aember),
	card.WithEnhance(card.Bonus.Damage),
	card.WithAbility(
		card.Trigger.Play, card.Exalt{
			Target: card.Target.Creature.Damaged(),
			Amount: 1,
		}),
)
