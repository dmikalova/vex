package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Finishing Blow
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Destroy a damaged creature -> steal 1 Æmber.
var FinishingBlow = card.New(
	"Finishing Blow",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 269),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Then{
			First:  card.Destroy{Target: card.Target.Creature.Damaged()},
			Result: card.StealAember{Amount: 1},
		}),
)
