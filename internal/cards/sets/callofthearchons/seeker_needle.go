package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Seeker Needle
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Common
//	Traits: Weapon
//
//	Action: Deal 1 damage to a Creature. If this damage destroys that Creature, gain 1 Æmber.
var SeekerNeedle = set.New(
	"Seeker Needle",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Common,
	card.Provenance(card.CotA, "290"),
	card.WithTraits(card.Traits.Weapon),
	card.WithAbility(
		card.Trigger.Action, card.DamageThen{
			Amount: 1,
			After:  card.IfDestroyed,
			Target: card.Target.Creature,
			Then: card.GainAember{
				Player: card.Controller,
				Amount: 1,
			},
		}),
)
