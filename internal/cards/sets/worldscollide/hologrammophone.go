package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Hologrammophone
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Common
//	Æmber:  1
//	Traits: Item
//
//	Action: Ward a creature.
var Hologrammophone = card.New(
	"Hologrammophone",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Common,
	card.Provenance(card.WC, "135"),
	card.WithAemberBonus(1),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.Action, card.Ward{Target: card.Target.Creature}),
)
