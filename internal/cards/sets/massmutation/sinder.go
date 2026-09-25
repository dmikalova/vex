package massmutation

import "github.com/dmikalova/vex/internal/card"

// Sinder
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  6
//	Armor:  2
//	Traits: Demon
//
//	Taunt.
//	Reap: Destroy a friendly creature.
var Sinder = set.New(
	"Sinder",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "013"),
	card.WithPower(6),
	card.WithArmor(2),
	card.WithTraits(card.Traits.Demon),
	card.WithKeywords(card.Keyword.Taunt),
	card.WithAbility(
		card.Trigger.Reap, card.Destroy{Target: card.Target.FriendlyCreature}),
)
