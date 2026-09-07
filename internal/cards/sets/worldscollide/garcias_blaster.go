package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Garcia's Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This creature gains, "Fight/Reap: You may choose one:
//	- Deal 2 damage to a creature
//	- Attach this creature to Sensor Chief Garcia, and steal 1 Æmber."
var GarciasBlaster = card.New(
	"Garcia's Blaster",
	card.House.StarAlliance,
	card.Type.Upgrade,
	// TODO(variant): rarity relabelled from Variant to Rare — handle manually
	card.Rarity.Rare,
	card.Provenance(card.WC, 347),
	card.WithAemberBonus(1),
	card.Connects(card.Pull(SensorChiefGarcia, 1)),
	card.WithStatic(card.StaticModifier{
		Granted: card.FightOrReap(card.May{Do: card.ChooseOne{Options: []card.Effect{
			card.DealDamage{Amount: 2, Target: card.Target.Creature},
			card.Sequence{Effects: []card.Effect{
				card.AttachSelfTo{Host: "Sensor Chief Garcia"},
				card.StealAember{Amount: 1},
			}},
		}}}),
	}),
)
