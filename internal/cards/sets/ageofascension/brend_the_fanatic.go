package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Brend the Fanatic
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Elf • Thief
//
//	Skirmish.
//	Play: Your opponent gains 1 Æmber.
//	Destroyed: Steal 3 Æmber.
var BrendTheFanatic = card.New(
	"Brend the Fanatic",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, 284),
	card.WithPower(3),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	card.WithKeywords(card.Keyword.Skirmish),
	card.WithAbility(
		card.Trigger.Play, card.GainAember{Player: card.Opponent, Amount: 1}),
	card.WithAbility(
		card.Trigger.Destroyed, card.StealAember{Amount: 3}),
)
