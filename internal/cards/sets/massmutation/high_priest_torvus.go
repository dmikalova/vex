package massmutation

import "github.com/dmikalova/vex/internal/card"

// High Priest Torvus
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Armor:  1
//	Traits: Dinosaur • Priest
//
//	Reap: You may exalt High Priest Torvus -> after you resolve your next tactic this turn, put it into your hand instead of your discard pile.
var HighPriestTorvus = set.New(
	"High Priest Torvus",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "222"),
	card.WithPower(4),
	card.WithArmor(1),
	card.WithTraits(card.Traits.Dinosaur, card.Traits.Priest),
	card.WithAbility(card.Trigger.Reap, card.May{Do: card.Then{
		First: card.Exalt{
			Target: card.Target.This,
			Amount: 1,
		},
		Result: card.PutNextTacticIntoHand{},
	}}),
)
