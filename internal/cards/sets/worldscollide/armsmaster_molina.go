package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Armsmaster Molina
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Human
//
//	Hazardous 3.
//	Each neighboring Creature gains hazardous 3.
var ArmsmasterMolina = card.New(
	"Armsmaster Molina",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "292"),
	card.InCluster(card.Pulled(molinasBlasterCluster, 1, 1.25)),
	card.WithPower(4),
	card.WithTraits(card.Traits.Human),
	card.WithHazardous(3),
	card.WithConstant(card.ConstantAbility{
		Target:         card.Target.EachCreature.Neighboring(),
		HazardousBonus: 3,
	}),
)
