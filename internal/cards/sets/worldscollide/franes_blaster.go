package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Frane's Blaster
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Special
//	Æmber:  1
//
//	This creature gains, "Fight/Reap: You may choose one:
//	- Deal 2 damage to a creature
//	- Attach this creature to First Officer Frane, and move all Æmber from First Officer Frane to your pool."
var FranesBlaster = card.New(
	"Frane's Blaster",
	card.House.StarAlliance,
	card.Type.Upgrade,
	// TODO(variant): rarity relabelled from Variant to Special — handle manually
	card.Rarity.Special,
	card.Provenance(card.WC, "346"),
	card.WithAemberBonus(1),
	card.Connects(card.Pull(FirstOfficerFrane, 1)),
	card.WithStatic(card.StaticModifier{
		Granted: card.FightOrReap(card.May{Do: card.ChooseOne{Options: []card.Effect{
			card.DealDamage{Amount: 2, Target: card.Target.Creature},
			card.Sequence{Effects: []card.Effect{
				card.AttachSelfTo{Host: "First Officer Frane"},
				card.MoveAember{
					All:  true,
					From: card.Target.AttachedHost.Named("First Officer Frane"),
					To:   card.Controller,
				},
			}},
		}}}),
	}),
)
