package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Gabos Longarms
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Traits: Demon
//
//	Before Fight: Choose a creature - Gabos Longarms deals its fight damage to the chosen creature instead of to the creature it is fighting.
var GabosLongarms = card.New(
	"Gabos Longarms",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 86),
	card.WithPower(5),
	card.WithTraits(card.Traits.Demon),
	card.WithAbility(
		card.Trigger.BeforeFight, card.RedirectFightDamage{Target: card.Target.Creature}),
)
