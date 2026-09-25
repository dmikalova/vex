package massmutation

import "github.com/dmikalova/vex/internal/card"

// Galeatops
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  12
//	Traits: Beast
//
//	Galeatops deals 4 damage when fighting.
var Galeatops = set.New(
	"Galeatops",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "193"),
	card.WithPower(12),
	card.WithTraits(card.Traits.Beast),
	card.WithAttackDamage(card.AttackDamage{
		Amount: 4,
		Fixed:  true,
	}),
)
