package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Helper Bot
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Robot
//
//	Play: Play a non-Logos card.
var HelperBot = card.New(
	"Helper Bot",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 112),
	card.WithPower(1),
	card.WithTraits(card.Traits.Robot),
	card.WithAbility(
		card.Trigger.Play, card.PlayFrom{
			From:   card.Hand,
			House:  card.House.Self,
			Except: true,
		}),
)
