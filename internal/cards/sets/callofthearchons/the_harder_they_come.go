package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// The Harder They Come
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Purge a creature with power 5 or higher.
var TheHarderTheyCome = card.New(
	"The Harder They Come",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 228),
	card.WithAbility(
		card.Trigger.Play, card.PurgeCreature{
			Target: card.Target.Creature.PowerAtLeast(5),
		}),
)
