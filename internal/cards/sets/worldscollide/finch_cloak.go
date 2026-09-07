package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Finch Cloak
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Elf • Thief
//
//	Fight/Reap: If your opponent has more Æmber than you, steal 1 Æmber. Otherwise, each player gains 1 Æmber.
var FinchCloak = card.New(
	"Finch Cloak",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, 267),
	card.WithPower(4),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	card.WithFightOrReap(card.Conditional{
		Cond: card.OpponentAember{Is: card.MoreThanYou},
		Then: card.StealAember{Amount: 1},
		Else: card.GainAember{
			Player: card.EachPlayer,
			Amount: 1,
		},
	}),
)
