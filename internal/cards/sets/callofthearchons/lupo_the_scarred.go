package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Lupo the Scarred
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  6
//	Traits: Beast
//
//	Skirmish.
//	Play: Deal 2 damage to an enemy creature.
var LupoTheScarred = card.New(
	"Lupo the Scarred",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 359),
	card.WithPower(6),
	card.WithTraits("Beast"),
	card.WithKeywords(card.Keyword.Skirmish),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount: 2,
			Target: card.Target.EnemyCreature,
		}),
)
