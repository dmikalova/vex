package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Inka the Spider
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Beast
//
//	Poison.
//	Play/Reap: Stun a creature.
var InkaTheSpider = card.New(
	"Inka the Spider",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 356),
	card.WithPower(1),
	card.WithTraits(card.Traits.Beast),
	card.WithKeywords(card.Keyword.Poison),
	card.WithPlayReap(card.Stun{Target: card.Target.Creature}),
)
