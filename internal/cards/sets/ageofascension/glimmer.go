package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Glimmer
//
//	House:  Untamed
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Faerie
//
//	Alpha.
//	Play: Put a card from your discard pile into your hand.
var Glimmer = card.New(
	"Glimmer",
	card.House.Untamed,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 323),
	card.WithPower(1),
	card.WithTraits(card.Traits.Faerie),
	card.WithKeywords(card.Keyword.Alpha),
	card.WithAbility(
		card.Trigger.Play, card.PutFromDiscard{Destination: card.To.Hand}),
)
