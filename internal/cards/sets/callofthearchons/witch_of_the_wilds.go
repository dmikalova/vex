package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Witch of the Wilds
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Beast • Witch
//
//	During each turn in which Untamed is not your active house, you may play one Untamed card.
var WitchOfTheWilds = card.New(
	"Witch of the Wilds",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 347),
	card.WithPower(4),
	card.WithTraits("Beast", "Witch"),
	card.WithPlayPermission(card.PlayPermission{House: card.House.Self, Count: 1}),
)
