package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Magda the Rat
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Elf • Thief
//
//	Elusive.
//	Play: Steal 2 Æmber.
//	Leaves Play: Your opponent steals 2 Æmber.
var MagdaTheRat = card.New(
	"Magda the Rat",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 303),
	card.WithPower(4),
	card.WithTraits("Elf", "Thief"),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(card.Trigger.Play, card.StealAember{Amount: 2}),
	card.WithAbility(
		card.Trigger.LeavesPlay, card.StealAember{
			Player: card.Opponent,
			Amount: 2,
		}),
)
