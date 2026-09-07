package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Way of the Bear
//
//	House:  Untamed
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	This creature gains +2 assault.
var WayOfTheBear = card.New(
	"Way of the Bear",
	card.House.Untamed,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "369"),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{AssaultBonus: 2}),
)
