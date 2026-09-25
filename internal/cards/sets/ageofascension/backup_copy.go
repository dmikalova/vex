package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Backup Copy
//
//	House:  Logos
//	Type:   Upgrade
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	This creature gains, "Destroyed: Put this creature on top of its owner's deck."
var BackupCopy = set.New(
	"Backup Copy",
	card.House.Logos,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "124"),
	card.WithBonus(card.Bonus.Aember),
	card.WithStatic(card.StaticModifier{
		Granted: []card.Ability{{
			Trigger: card.Trigger.Destroyed,
			Effect: card.PutFromPlay{
				Target:      card.Target.This,
				Destination: card.To.TopOfDeck,
			},
		}},
	}),
)
