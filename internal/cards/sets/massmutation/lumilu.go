package massmutation

import "github.com/dmikalova/vex/internal/card"

// Lumilu
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  4
//	Traits: Beast • Cat
//
//	Reap: For each other friendly Beast creature, gain 1 Æmber.
var Lumilu = set.New(
	"Lumilu",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "385"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Beast, card.Traits.Cat),
	card.WithAbility(
		card.Trigger.Reap, card.GainAember{
			Player: card.Controller,
			Amount: 1,
			Per: card.CardsInPlay{
				Player: card.Controller,
				Type:   card.Type.Creature,
				Trait:  card.Traits.Beast,
				Other:  true,
			},
		}),
)
