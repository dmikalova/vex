package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Faygin
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Human • Thief
//
//	Elusive.
//	Reap: Put an Urchin from play or from your discard pile into your hand.
var Faygin = card.New(
	"Faygin",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 300),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human, card.Traits.Thief),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Reap, card.ReturnNamedToHand{Name: "Urchin"}),
)
