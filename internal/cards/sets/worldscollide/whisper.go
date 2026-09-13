package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Whisper
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Elf • Thief
//
//	Elusive.
//	Action: Lose 1 Æmber -> destroy a Creature.
var Whisper = card.New(
	"Whisper",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "265"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Elf, card.Traits.Thief),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Action, card.Then{
			First: card.LoseAember{
				Player: card.Controller,
				Amount: 1,
			},
			Result: card.Destroy{Target: card.Target.Creature},
		}),
)
