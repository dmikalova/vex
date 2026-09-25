package massmutation

import "github.com/dmikalova/vex/internal/card"

// Bonesaw
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Demon
//
//	If a friendly creature has been destroyed this turn, Bonesaw enters play ready.
var Bonesaw = set.New(
	"Bonesaw",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "002"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Demon),
	card.WithEntersPlay(card.Conditional{
		Cond: card.CreatureDestroyedThisTurn{Player: card.Controller},
		Then: card.Ready{Target: card.Target.This},
	}),
)
