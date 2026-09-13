package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Wild Spirit
//
//	House:  Untamed
//	Type:   Upgrade
//	Rarity: Uncommon
//	Æmber:  1
//
//	This Creature gains, "Reap: this Creature captures 1 Æmber from your opponent."
var WildSpirit = card.New(
	"Wild Spirit",
	card.House.Untamed,
	card.Type.Upgrade,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "384"),
	card.WithAemberBonus(1),
	card.WithStatic(card.StaticModifier{
		Granted: []card.Ability{{
			Trigger: card.Trigger.Reap,
			Effect: card.CaptureAember{
				Amount: 1,
				Target: card.Target.This,
				Source: card.Opponent,
			},
		}},
	}),
)
