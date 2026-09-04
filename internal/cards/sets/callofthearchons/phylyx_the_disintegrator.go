package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Phylyx the Disintegrator
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Martian • Soldier
//
//	Elusive.
//	Action: For each other friendly Mars creature, your opponent loses 1 Æmber.
var PhylyxTheDisintegrator = card.New("Phylyx the Disintegrator",
	card.House.Mars, card.Type.Creature, card.Rarity.Rare,
	card.Provenance(card.CotA, 197),
	card.WithPower(1),
	card.WithTraits(card.Traits.Martian, card.Traits.Soldier),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Action, card.LoseAember{
			Player: card.Opponent,
			Amount: 1,
			Per: card.InPlay{
				Player: card.Controller,
				Type:   card.Type.Creature,
				House:  card.House.Self,
				Other:  true,
			},
		}),
)
