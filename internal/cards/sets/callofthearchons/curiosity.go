package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Curiosity
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Destroy each Scientist creature.
var Curiosity = set.New(
	"Curiosity",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "320"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play,
		card.Destroy{Target: card.Target.EachCreature.WithTrait(card.Traits.Scientist)},
	),
)
