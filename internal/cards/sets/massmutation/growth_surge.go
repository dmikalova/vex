package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Growth Surge
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Choose a flank creature. Give it three +1 power counters, its neighbor two +1 power counters, and the neighbor's other neighbor a +1 power counter.
var GrowthSurge = set.New(
	"Growth Surge",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "383"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.AddPowerCounter{Walk: []int{3, 2, 1}}),
)
