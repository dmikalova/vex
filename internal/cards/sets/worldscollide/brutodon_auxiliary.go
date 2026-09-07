package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Brutodon Auxiliary
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  6
//	Traits: Beast
//
//	Taunt, Hazardous 2.
var BrutodonAuxiliary = card.New(
	"Brutodon Auxiliary",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "183"),
	card.WithPower(6),
	card.WithTraits(card.Traits.Beast),
	card.WithKeywords(card.Keyword.Taunt),
	card.WithHazardous(2),
)
