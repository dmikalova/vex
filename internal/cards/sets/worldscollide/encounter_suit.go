package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Encounter Suit
//
//	House:  Star Alliance
//	Type:   Upgrade
//	Rarity: Rare
//
//	This creature gains, "After a Tactic is played but before it resolves, ward this creature."
var EncounterSuit = card.New(
	"Encounter Suit",
	card.House.StarAlliance,
	card.Type.Upgrade,
	card.Rarity.Rare,
	card.Provenance(card.WC, "330"),
	card.WithStatic(card.StaticModifier{
		Granted: []card.Ability{{
			Trigger: card.Trigger.AfterActionPlayedBeforeResolve,
			Effect:  card.Ward{Target: card.Target.This},
		}},
	}),
)
