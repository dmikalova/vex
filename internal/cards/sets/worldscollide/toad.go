package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Toad
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Connected
//	Power:  1
//	Traits: Beast
//
//	Toad cannot reap.
var Toad = set.New(
	"Toad",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Connected,
	card.Provenance(card.WC, "405"),
	card.InCluster(card.Pulled(xenosBloodshadowCluster, 1, 1)),
	card.WithPower(1),
	card.WithTraits(card.Traits.Beast),
	card.WithCannotBeUsedTo(card.UseKind.Reap),
)
