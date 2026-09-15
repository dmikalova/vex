package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Gizelhart's Zealot
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Human • Knight
//
//	Gizelhart's Zealot enters play ready and enrage Gizelhart's Zealot.
var GizelhartsZealot = set.New(
	"Gizelhart's Zealot",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "136"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Human, card.Traits.Knight),
	card.WithEntersPlay(card.Sequence{Effects: []card.Effect{
		card.Ready{Target: card.Target.This},
		card.Enrage{Target: card.Target.This},
	}}),
)
