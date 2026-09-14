package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Draining Touch
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Destroy a Creature with no Æmber on it.
var DrainingTouch = set.New(
	"Draining Touch",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "72"),
	card.WithAbility(
		card.Trigger.Play, card.Destroy{
			Target: card.Target.Creature.WithoutAember(),
		}),
)
