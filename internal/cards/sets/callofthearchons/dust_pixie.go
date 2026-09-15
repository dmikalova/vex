package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Dust Pixie
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Bonus:  Æmber Æmber
//	Traits: Faerie
var DustPixie = set.New(
	"Dust Pixie",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, "351"),
	card.WithBonus(card.Bonus.Aember, card.Bonus.Aember),
	card.WithPower(1),
	card.WithTraits(card.Traits.Faerie),
)
