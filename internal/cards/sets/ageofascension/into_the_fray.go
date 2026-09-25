package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Into the Fray
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//
//	Play: A friendly Brobnar creature gains, "Fight: Ready this creature."
var IntoTheFray = set.New(
	"Into the Fray",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, "13"),
	card.WithAbility(
		card.Trigger.Play, card.GainAbility{
			Target: card.Target.FriendlyCreature.House(card.Houses.Named(card.House.Self)),
			Ability: card.Ability{
				Trigger: card.Trigger.Fight,
				Effect:  card.Ready{Target: card.Target.This},
			},
		}),
)
