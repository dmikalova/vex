package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Bellowing Patrizate
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  7
//	Traits: Giant
//
//	After a creature enters play, if Bellowing Patrizate is ready, deal 1 damage to it.
var BellowingPatrizate = card.New(
	"Bellowing Patrizate",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 34),
	card.WithPower(7),
	card.WithTraits(card.Traits.Giant),
	card.WithAbility(
		card.Trigger.AfterCreatureEnters, card.Conditional{
			Cond: card.SourceReady{},
			Then: card.DealDamage{
				Amount: 1,
				Target: card.Target.Triggering,
			},
		}),
)
