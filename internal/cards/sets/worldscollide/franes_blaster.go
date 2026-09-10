package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Frane's Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Rare
//	Æmber:  1
//
//	This creature gains, "Fight/Reap: Choose one:
//	- Deal 2 damage to a creature
//	- Attach Frane's Blaster to First Officer Frane -> move all Æmber from First Officer Frane to your pool."
var FranesBlaster = card.New(
	"Frane's Blaster",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, "346"),
	card.Connects(card.Pull(FirstOfficerFrane, 1)),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		Granted: card.FightOrReap(card.ChooseOne{Options: []card.Effect{
			card.DealDamage{Amount: 2, Target: card.Target.Creature},
			card.Then{
				First: card.AttachSelfTo{Host: FirstOfficerFrane.Name},
				Result: card.MoveAember{
					All:  true,
					From: card.Target.AttachedHost.Named(FirstOfficerFrane.Name),
					To:   card.Controller,
				},
			},
		}}),
	}),
)
