package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Duskrunner
//
//	House:  Shadows
//	Type:   Upgrade
//	Rarity: Uncommon
//
//	This Creature gains, "Reap: Steal 1 Æmber."
var Duskrunner = card.New(
	"Duskrunner",
	card.House.Shadows,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "316"),
	card.WithStatic(card.StaticModifier{
		Granted: []card.Ability{
			{Trigger: card.Trigger.Reap, Effect: card.StealAember{Amount: 1}},
		},
	}),
)
