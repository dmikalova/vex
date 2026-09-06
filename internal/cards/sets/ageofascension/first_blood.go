package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// First Blood
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Alpha.
//	Play: Deal 2 damage for each friendly Brobnar creature, divided among any number of creatures.
var FirstBlood = card.New(
	"First Blood",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.AoA, 7),
	card.WithAemberBonus(1),
	card.WithKeywords(card.Keyword.Alpha),
	card.WithAbility(
		card.Trigger.Play, card.DealDamage{
			Spread: card.DivideDamage{
				Amount: 2,
				Per: card.InPlay{
					Player: card.Controller,
					Type:   card.Type.Creature,
					House:  card.House.Self,
				},
			},
		}),
)
