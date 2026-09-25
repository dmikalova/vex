package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Collar of Subordination
//
//	House:  Dis
//	Type:   Upgrade
//	Rarity: Rare
//
//	Play: Take control of this creature until Collar of Subordination leaves play.
var CollarOfSubordination = set.New(
	"Collar of Subordination",
	card.House.Dis,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "105"),
	card.WithAbility(
		card.Trigger.Play, card.TakeControl{
			Duration: card.Duration.UntilThisLeavesPlay,
		}),
)
