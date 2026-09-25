package massmutation

import "github.com/dmikalova/vex/internal/card"

// Gloriana's Attendant
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Mutant
//
//	Enhance Æmber Æmber.
var GlorianasAttendant = set.New(
	"Gloriana's Attendant",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "368"),
	card.WithEnhance(card.Bonus.Aember, card.Bonus.Aember),
	card.WithPower(1),
	card.WithTraits(card.Traits.Mutant),
)
