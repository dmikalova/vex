package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Champion Anaphiel
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Common
//	Power:  6
//	Armor:  1
//	Traits: Knight • Spirit
//
//	Taunt.
var ChampionAnaphiel = card.New(
	"Champion Anaphiel",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 239),
	card.WithPower(6),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Knight, card.Traits.Spirit),
	card.WithKeywords(card.Keyword.Taunt),
)
