package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Oubliette
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Purge a creature with power 3 or lower.
var Oubliette = card.New(
	"Oubliette",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 278),
	card.WithAbility(
		card.Trigger.Play, card.PurgeCreature{
			Target: card.Target.Creature.PowerAtMost(3),
		}),
)
