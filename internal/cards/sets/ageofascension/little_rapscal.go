package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Little Rapscal
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Goblin
//
//	Elusive.
//	Creatures must fight when used, if able.
var LittleRapscal = card.New(
	"Little Rapscal",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "25"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Goblin),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithRestrictions(card.Restrictions{MustFightIfAble: true}),
)
