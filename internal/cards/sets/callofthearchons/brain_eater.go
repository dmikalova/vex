package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Brain Eater
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  6
//	Traits: Cyborg • Beast
//
//	After a creature is destroyed in a fight with Brain Eater, draw a card.
var BrainEater = set.New(
	"Brain Eater",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "137"),
	card.WithPower(6),
	card.WithTraits(card.Traits.Cyborg, card.Traits.Beast),
	card.WithAbility(card.Trigger.AfterDestroyedFighting, card.Draw{Amount: 1}),
)
