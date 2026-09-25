package massmutation

import "github.com/dmikalova/vex/internal/card"

// Rad Penny
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Mutant • Thief
//
//	Play: Steal 1 Æmber.
//	Destroyed: Shuffle Rad Penny into its owner's deck.
var RadPenny = set.New(
	"Rad Penny",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "255"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Mutant, card.Traits.Thief),
	card.WithAbility(
		card.Trigger.Play, card.StealAember{Amount: 1}),
	card.WithAbility(
		card.Trigger.Destroyed, card.PutFromPlay{
			Target:      card.Target.This,
			Destination: card.To.DeckShuffled,
		}),
)
