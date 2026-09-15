package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Keyfrog
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Beast
//
//	Destroyed: Forge a key at current cost -> purge Keyfrog.
var Keyfrog = set.New(
	"Keyfrog",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "369"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Beast),
	card.WithAbility(
		card.Trigger.Destroyed, card.ForgeKey{}),
)
