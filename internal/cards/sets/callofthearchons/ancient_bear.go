package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Ancient Bear
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Beast
//
//	Assault 2.
var AncientBear = card.New(
	"Ancient Bear",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 345),
	card.WithPower(5),
	card.WithTraits(card.Traits.Beast),
	card.WithAssault(2),
)
