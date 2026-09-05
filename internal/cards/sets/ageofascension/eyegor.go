package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Eyegor
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Cyborg
//
//	Play: Look at the top 3 cards of your deck, put 1 into your hand, and discard the others.
var Eyegor = card.New(
	"Eyegor",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 111),
	card.WithPower(2),
	card.WithTraits(card.Traits.Cyborg),
	card.WithAbility(
		card.Trigger.Play, card.LookAtTop{Amount: 3}),
)
