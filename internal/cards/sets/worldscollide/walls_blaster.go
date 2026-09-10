package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Walls' Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This creature gains, "Fight/Reap: Choose one:
//	- Deal 2 damage to a creature
//	- Attach Walls' Blaster to Chief Engineer Walls -> for each upgrade on Chief Engineer Walls, stun a creature."
var WallsBlaster = card.New(
	"Walls' Blaster",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, "352"),
	card.Connects(card.Pull(ChiefEngineerWalls, 1)),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		Granted: card.FightReap(card.ChooseOne{Options: []card.Effect{
			card.DealDamage{Amount: 2, Target: card.Target.Creature},
			card.Then{
				First: card.AttachSelfTo{Host: ChiefEngineerWalls.Name},
				Result: card.Repeat{
					Times: card.UpgradesOn{
						Target: card.Target.AttachedHost.Named(ChiefEngineerWalls.Name),
					},
					Do: card.Stun{Target: card.Target.Creature},
				},
			},
		}}),
	}),
)
