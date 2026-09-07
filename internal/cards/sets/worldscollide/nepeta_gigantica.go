package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Nepeta Gigantica
//
//	House:  Untamed
//	Type:   Artifact
//	Rarity: Rare
//	Traits: Item
//
//	Action: Choose one:
//	- Stun a creature with power 5 or higher
//	- Stun a Giant creature.
var NepetaGigantica = card.New(
	"Nepeta Gigantica",
	card.House.Untamed,
	card.Type.Artifact,
	card.Rarity.Rare,
	card.Provenance(card.WC, "394"),
	card.WithTraits(card.Traits.Item),
	card.WithAbility(
		card.Trigger.Action, card.ChooseOne{Options: []card.Effect{
			card.Stun{Target: card.Target.Creature.PowerAtLeast(5)},
			card.Stun{Target: card.Target.Creature.WithTrait(card.Traits.Giant)},
		}}),
)
