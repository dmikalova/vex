package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Garcia's Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This creature gains, "Fight/Reap: Choose one:
//	- Deal 2 damage to a creature
//	- Attach Garcia's Blaster to Sensor Chief Garcia -> steal 1 Æmber."
var GarciasBlaster = card.New(
	"Garcia's Blaster",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, "347"),
	card.Connects(card.Pull(SensorChiefGarcia, 1)),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		Granted: card.FightReap(card.ChooseOne{Options: []card.Effect{
			card.DealDamage{Amount: 2, Target: card.Target.Creature},
			card.Then{
				First:  card.AttachSelfTo{Host: SensorChiefGarcia.Name},
				Result: card.StealAember{Amount: 1},
			},
		}}),
	}),
)
