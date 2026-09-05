package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Drummernaut
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  6
//	Traits: Giant
//
//	Play/Fight/Reap: Put another friendly Giant trait creature into its owner's hand.
var Drummernaut = card.New(
	"Drummernaut",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, 6),
	card.WithPower(6),
	card.WithTraits(card.Traits.Giant),
	card.WithPlayFightReap(card.PutFromPlay{
		Target:      card.Target.OtherFriendlyCreature.WithTrait(card.Traits.Giant),
		Destination: card.To.Hand,
	}),
)
