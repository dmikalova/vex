package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Way of the Bear
//
//	House:  Untamed
//	Type:   Upgrade
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	This creature gains +2 assault.
var WayOfTheBear = set.New(
	"Way of the Bear",
	card.House.Untamed,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "369"),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{AssaultBonus: 2}),
)
