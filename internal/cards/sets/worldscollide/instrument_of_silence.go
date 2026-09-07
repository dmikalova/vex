package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Instrument of Silence
//
//	House:  Untamed
//	Type:   Upgrade
//	Rarity: Uncommon
//
//	This creature gains skirmish.
//	This creature gains, "Fight: Gain 1 Æmber."
var InstrumentOfSilence = card.New(
	"Instrument of Silence",
	card.House.Untamed,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "375"),
	card.WithStatic(card.StaticModifier{
		Keywords: card.Keywords(card.Keyword.Skirmish),
		Granted: []card.Ability{{
			Trigger: card.Trigger.Fight,
			Effect: card.GainAember{
				Player: card.Controller,
				Amount: 1,
			},
		}},
	}),
)
