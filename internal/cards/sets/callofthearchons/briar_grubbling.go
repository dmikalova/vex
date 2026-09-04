package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Briar Grubbling
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Beast • Insect
//
//	Hazardous 5.
var BriarGrubbling = card.New(
	"Briar Grubbling",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 348),
	card.WithPower(2),
	card.WithTraits(card.Traits.Beast, card.Traits.Insect),
	card.WithHazardous(5),
)
