package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Drummernaut
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  6
//	Traits: Giant
//
//	Play/Fight/Reap: Put another friendly Giant creature into its owner's hand.
var Drummernaut = set.New(
	"Drummernaut",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, "6"),
	card.WithPower(6),
	card.WithTraits(card.Traits.Giant),
	card.WithAbility(card.Trigger.PlayFightReap, card.PutFromPlay{
		Target:      card.Target.OtherFriendlyCreature.WithTrait(card.Traits.Giant),
		Destination: card.To.Hand,
	}),
)
