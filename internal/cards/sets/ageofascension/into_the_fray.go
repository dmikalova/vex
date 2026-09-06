package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Into the Fray
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//
//	Play: A friendly Brobnar creature gains, "Fight: Ready this creature."
var IntoTheFray = card.New(
	"Into the Fray",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, 13),
	card.WithAbility(
		card.Trigger.Play, card.GainAbility{
			Target: card.Target.FriendlyCreature.OfHouse(card.House.Self),
			Ability: card.Ability{
				Trigger: card.Trigger.Fight,
				Effect:  card.Ready{Target: card.Target.This},
			},
		}),
)
