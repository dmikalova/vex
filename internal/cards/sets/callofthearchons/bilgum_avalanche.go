package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Bilgum Avalanche
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  5
//	Traits: Giant
//
//	After you forge a key, deal 2 damage to each enemy creature.
var BilgumAvalanche = card.New(
	"Bilgum Avalanche",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 28),
	card.WithPower(5),
	card.WithTraits(card.Traits.Giant),
	card.WithAbility(
		card.Trigger.AfterForgeKey, card.DealDamage{
			Amount: 2,
			Target: card.Target.EachEnemyCreature,
		}),
)
