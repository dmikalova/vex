package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Doc Bookton
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Human • Scientist
//
//	Reap: Draw a card.
var DocBookton = card.New(
	"Doc Bookton",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 139),
	card.WithPower(5),
	card.WithTraits(card.Traits.Human, card.Traits.Scientist),
	card.WithAbility(
		card.Trigger.Reap, card.Draw{Amount: 1}),
)
