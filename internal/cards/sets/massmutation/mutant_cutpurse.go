package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Mutant Cutpurse
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Mutant • Thief
//
//	Enhance Damage Damage Damage.
var MutantCutpurse = set.New(
	"Mutant Cutpurse",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "252"),
	card.WithEnhance(card.Bonus.Damage, card.Bonus.Damage, card.Bonus.Damage),
	card.WithPower(3),
	card.WithTraits(card.Traits.Mutant, card.Traits.Thief),
)
