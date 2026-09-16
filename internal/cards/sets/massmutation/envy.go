package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Envy
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Special
//	Power:  3
//	Traits: Demon • Sin
//
//	Elusive.
//	Reap: If there are 2 or more friendly Sin creatures in play, Envy captures all your opponent's Æmber.
var Envy = set.New(
	"Envy",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Special,
	card.Provenance(card.MM, "056"),
	card.InCluster(sinsCluster),
	card.OneCopyPerDeck(),
	card.WithPower(3),
	card.WithTraits(card.Traits.Demon, card.Traits.Sin),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(card.Trigger.Reap, card.Conditional{
		Cond: card.InPlay{
			Player: card.Controller,
			Type:   card.Type.Creature,
			Trait:  card.Traits.Sin,
			Amount: 2,
		},
		Then: card.CaptureAember{
			All:    true,
			Target: card.Target.This,
			Source: card.Opponent,
		},
	}),
)
