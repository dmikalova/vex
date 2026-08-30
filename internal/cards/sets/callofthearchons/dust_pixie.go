package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Dust Pixie
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Æmber:  2
//	Traits: Faerie
var DustPixie = card.New(
	"Dust Pixie",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 351),
	card.WithPower(1),
	card.WithAemberBonus(2),
	card.WithTraits("Faerie"),
)
