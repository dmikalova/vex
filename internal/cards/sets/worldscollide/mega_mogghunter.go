package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Mega Mogghunter
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Connected
//	Power:  8
//	Traits: Giant
//
//	Fight: Deal 2 damage to a flank Creature.
var MegaMogghunter = set.New(
	"Mega Mogghunter",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Connected,
	card.Provenance(card.WC, "59"),
	card.InCluster(card.Pulled(mogghuntersBrewCluster, 1, 1.25)),
	card.WithPower(8),
	card.WithTraits(card.Traits.Giant),
	card.WithAbility(
		card.Trigger.Fight, card.DealDamage{
			Target: card.Target.Creature.OnFlank(),
			Amount: 2,
		}),
)
