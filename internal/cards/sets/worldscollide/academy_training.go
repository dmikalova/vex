package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// AcademyTraining
//
//	House:  Logos
//	Type:   Upgrade
//	Rarity: Rare
//
//	If you control this creature, it belongs to house Logos. (Instead of its original house.)
//	This creature gains, "Reap: Draw a card."
var AcademyTraining = card.New(
	"Academy Training",
	card.House.Logos,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, "161"),
	card.WithStatic(card.StaticModifier{
		HouseOverride: card.House.Self,
		Granted: []card.Ability{{
			Trigger: card.Trigger.Reap,
			Effect:  card.Draw{Amount: 1},
		}},
	}),
)
