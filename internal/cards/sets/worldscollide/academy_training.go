package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Academy Training
//
//	House:  Logos
//	Type:   Upgrade
//	Rarity: Rare
//
//	This creature belongs to Logos and this creature gains "Reap: Draw a card."
var AcademyTraining = set.New(
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
