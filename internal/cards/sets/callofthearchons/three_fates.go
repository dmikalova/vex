package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Three Fates
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Destroy the 3 most powerful creatures.
var ThreeFates = set.New(
	"Three Fates",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "71"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play,
		card.Destroy{Target: card.Target.EachCreature.Refine(card.MostPowerfulN(3))},
	),
)
