package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// The Spirit's Way
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Destroy each creature with power 3 or higher.
var TheSpiritsWay = card.New(
	"The Spirit's Way",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 229),
	card.WithAbility(
		card.Trigger.Play, card.Destroy{
			Target: card.Target.EachCreature.PowerAtLeast(3),
		}),
)
