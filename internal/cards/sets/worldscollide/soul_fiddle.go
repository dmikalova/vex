package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Soul Fiddle
//
//	House:  Dis
//	Type:   Artifact
//	Rarity: Uncommon
//	Traits: Item
//
//	Action: Enrage a creature.
var SoulFiddle = card.New(
	"Soul Fiddle",
	card.House.Dis,
	card.Type.Artifact,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "97"),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.Action, card.Enrage{Target: card.Target.Creature}),
)
