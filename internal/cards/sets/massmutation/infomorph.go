package massmutation

import "github.com/dmikalova/vex/internal/card"

// Infomorph
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Mutant
//
//	Enhance Draw Draw.
var Infomorph = set.New(
	"Infomorph",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "074"),
	card.WithEnhance(card.Bonus.Draw, card.Bonus.Draw),
	card.WithPower(4),
	card.WithTraits(card.Traits.Mutant),
)
