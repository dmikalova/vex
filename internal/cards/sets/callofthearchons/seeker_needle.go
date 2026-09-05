package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Seeker Needle
//
//	House:  Shadows
//	Type:   Artifact
//	Rarity: Common
//	Traits: Weapon
//
//	Action: Deal 1 damage to a creature. If this damage destroys that creature, gain 1 Æmber.
var SeekerNeedle = card.New(
	"Seeker Needle",
	card.House.Shadows,
	card.Type.Artifact,
	card.Rarity.Common,
	card.Provenance(card.CotA, 290),
	card.WithTraits(card.Traits.Weapon),
	card.WithAbility(
		card.Trigger.Action, card.DamageThenIfDestroyed{
			Amount: 1,
			Target: card.Target.Creature,
			Then: card.GainAember{
				Player: card.Controller,
				Amount: 1,
			},
		}),
)
