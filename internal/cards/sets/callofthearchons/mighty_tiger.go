package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Mighty Tiger
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Beast
//
//	Play: Deal 4 damage to an enemy creature.
var MightyTiger = card.New(
	"Mighty Tiger",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 360),
	card.WithPower(4),
	card.WithTraits("Beast"),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount: 4,
			Target: card.Target.EnemyCreature,
		}),
)
