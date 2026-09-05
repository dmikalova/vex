package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Skybooster Squadron
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Martian • Soldier
//
//	Reap: Put Skybooster Squadron into its owner's hand.
var SkyboosterSquadron = card.New(
	"Skybooster Squadron",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 170),
	card.WithPower(4),
	card.WithTraits(card.Traits.Martian, card.Traits.Soldier),
	card.WithAbility(
		card.Trigger.Reap, card.PutFromPlay{
			Target:      card.Target.This,
			Destination: card.To.Hand,
		}),
)
