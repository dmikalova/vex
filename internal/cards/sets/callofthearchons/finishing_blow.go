package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Finishing Blow
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Destroy a damaged creature -> steal 1 Æmber.
var FinishingBlow = set.New(
	"Finishing Blow",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "269"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Then{
			First:  card.Destroy{Target: card.Target.Creature.Damaged()},
			Result: card.StealAember{Amount: 1},
		}),
)
