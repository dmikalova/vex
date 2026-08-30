package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Biomatrix Backup
//
//	House:  Mars
//	Type:   Upgrade
//	Rarity: Common
//	Æmber:  1
//
//	This creature gains, "Destroyed: Put this creature into its owner's archives."
var BiomatrixBackup = card.New(
	"Biomatrix Backup",
	card.House.Mars,
	card.Type.Upgrade,
	card.Rarity.Common,
	card.Provenance(card.CotA, 208),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		Granted: []card.Ability{
			{Trigger: card.Trigger.Destroyed, Effect: card.PutFromPlay{
				Target:      card.Target.This,
				Destination: card.To.Archives,
			}},
		},
	}),
)
