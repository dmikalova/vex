package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Hologrammophone
//
//	House:  Logos
//	Type:   Artifact
//	Rarity: Common
//	Bonus:  Æmber
//	Traits: Item
//
//	Action: Ward a creature.
var Hologrammophone = set.New(
	"Hologrammophone",
	card.House.Logos,
	card.Type.Artifact,
	card.Rarity.Common,
	card.Provenance(card.WC, "135"),
	card.WithBonus(card.Bonus.Aember),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.Action, card.Ward{Target: card.Target.Creature}),
)
