package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Tolas
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  1
//	Traits: Imp
//
//	Elusive.
//	Each creature gains, "Destroyed: Your opponent gains 1 Æmber."
var Tolas = card.New(
	"Tolas",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 103),
	card.WithPower(1),
	card.WithTraits("Imp"),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithConstantAbility(card.ConstantAbility{
		Target: card.Target.EachCreature,
		Granted: []card.Ability{{
			Trigger: card.Trigger.Destroyed,
			Effect: card.GainAember{
				Player: card.Opponent,
				Amount: 1,
			},
		}},
	}),
)
