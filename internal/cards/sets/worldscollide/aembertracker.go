package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Aembertracker
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Beast
//
//	Play: Deal 2 damage to each enemy Creature with Æmber on it, ignoring armor.
var Aembertracker = set.New(
	"Aembertracker",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "324"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Beast),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Amount:      2,
			Target:      card.Target.EachEnemyCreature.WithAember(),
			IgnoreArmor: true,
		}),
)
