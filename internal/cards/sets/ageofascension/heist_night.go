package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Heist Night
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//
//	Alpha.
//	Play: For each friendly Thief trait creature, steal 1 Æmber.
var HeistNight = card.New(
	"Heist Night",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, 303),
	card.WithKeywords(card.Keyword.Alpha),
	card.WithAbility(
		card.Trigger.Play, card.StealAember{
			Amount: 1,
			Per: card.InPlay{
				Player: card.Controller,
				Type:   card.Type.Creature,
				Trait:  card.Traits.Thief,
			},
		}),
)
