package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Kindrith Longshot
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Human • Ranger
//
//	Elusive, Skirmish.
//	Reap: Deal 2 damage to a creature.
var KindrithLongshot = card.New(
	"Kindrith Longshot",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 357),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human, card.Traits.Ranger),
	card.WithKeywords(card.Keyword.Elusive, card.Keyword.Skirmish),
	card.WithAbility(
		card.Trigger.Reap, card.DealDamage{
			Amount: 2,
			Target: card.Target.Creature,
		}),
)
