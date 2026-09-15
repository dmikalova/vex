package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Liam Say
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Elf
//
//	Elusive.
//	At the start of your turn, you may deal 1 damage to a creature.
var LiamSay = set.New(
	"Liam Say",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "284"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Elf),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.StartOfTurn, card.May{
			Do: card.DealDamage{
				Amount: 1,
				Target: card.Target.Creature,
			},
		}),
)
