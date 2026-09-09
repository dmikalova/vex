package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Toad
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Special
//	Power:  1
//	Traits: Beast
//
//	Toad cannot reap.
var Toad = card.New(
	"Toad",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Special,
	card.Provenance(card.WC, "405"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Beast),
	card.WithCannotBeUsedTo(card.UseKind.Reap),
)
