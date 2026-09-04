package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Niffle Ape
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Beast • Niffle
//
//	While Niffle Ape is attacking, ignore taunt and elusive.
var NiffleApe = card.New(
	"Niffle Ape",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 363),
	card.WithPower(3),
	card.WithTraits(card.Traits.Beast, card.Traits.Niffle),
	card.WithAttackIgnores(card.Keyword.Taunt, card.Keyword.Elusive),
)
