package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Pingle Who Annoys
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Rare
//	Power:  2
//	Traits: Goblin
//
//	Elusive.
//	Play: Deal 1 damage to each enemy creature.
var PingleWhoAnnoys = card.New(
	"Pingle Who Annoys",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 43),
	card.WithPower(2),
	card.WithTraits("Goblin"),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount: 1,
			Target: card.Target.EachEnemyCreature,
		}),
)
