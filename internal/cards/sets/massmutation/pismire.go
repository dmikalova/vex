package massmutation

import "github.com/dmikalova/vex/internal/card"

// Pismire
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Mutant
//
//	While you control more Mutant creatures than your opponent, your opponent's keys cost +2 Æmber.
var Pismire = set.New(
	"Pismire",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "372"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Mutant),
	card.WithKeyCost(card.KeyCostChange(card.Opponent, 2).While(
		card.ControlsMoreCreatures{Trait: card.Traits.Mutant},
	)),
)
