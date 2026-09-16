package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Splinter
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Mutant • Thief
//
//	Enhance Damage Damage Damage Damage Damage Damage.
var Splinter = set.New(
	"Splinter",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "294"),
	card.WithEnhance(
		card.Bonus.Damage,
		card.Bonus.Damage,
		card.Bonus.Damage,
		card.Bonus.Damage,
		card.Bonus.Damage,
		card.Bonus.Damage,
	),
	card.WithPower(1),
	card.WithTraits(card.Traits.Mutant, card.Traits.Thief),
)
