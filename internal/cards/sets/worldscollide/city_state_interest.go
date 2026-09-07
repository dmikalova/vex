package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// City-State Interest
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Each friendly creature captures 1 Æmber from your own side.
var CityStateInterest = card.New(
	"City-State Interest",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 200),
	card.WithAbility(
		card.Trigger.Play, card.CaptureAember{
			Amount: 1,
			Target: card.Target.EachFriendlyCreature,
			Source: card.Controller,
		}),
)
