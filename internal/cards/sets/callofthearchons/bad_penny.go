package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Bad Penny
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Human • Thief
//
//	Destroyed: Put Bad Penny into its owner's hand.
var BadPenny = card.New(
	"Bad Penny",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 296),
	card.WithPower(1),
	card.WithTraits("Human", "Thief"),
	card.WithAbility(
		card.Trigger.Destroyed, card.PutFromPlay{
			Target:      card.Target.This,
			Destination: card.To.Hand,
		}),
)
