package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Bo Nithing
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Elf • Thief
//
//	Play: For each forged key your opponent has, steal 1 Æmber.
var BoNithing = set.New(
	"Bo Nithing",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "245"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	card.WithAbility(
		card.Trigger.Play, card.StealAember{
			Amount: 1,
			Per:    card.ForgedKeys{Player: card.Opponent},
		}),
)
