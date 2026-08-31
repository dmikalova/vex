package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Three Fates
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Destroy the 3 most powerful creatures.
var ThreeFates = card.New(
	"Three Fates",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 71),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play,
		card.Destroy{Target: card.Target.EachCreature.Selector(card.MostPowerful(3))},
	),
)
