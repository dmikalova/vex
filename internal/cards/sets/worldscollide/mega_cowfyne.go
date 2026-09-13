package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Mega Cowfyne
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Special
//	Power:  7
//	Traits: Giant
//
//	Before Fight: Deal 2 damage to each neighbor of the Creature Mega Cowfyne fights.
var MegaCowfyne = card.New(
	"Mega Cowfyne",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Special,
	card.Provenance(card.WC, "55"),
	card.WithPower(7),
	card.WithTraits(card.Traits.Giant),
	card.WithAbility(
		card.Trigger.BeforeFight, card.DealDamage{
			Amount: 2,
			Target: card.Target.CreatureFought.NeighborsOf(),
		}),
)
