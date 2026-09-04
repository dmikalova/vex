package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Lady Maxena
//
//	House:  Sanctum
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  5
//	Traits: Knight • Spirit
//
//	Play: Stun a creature.
//	Action: Put Lady Maxena into its owner's hand.
var LadyMaxena = card.New(
	"Lady Maxena",
	card.House.Sanctum,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 251),
	card.WithPower(5),
	card.WithTraits(card.Traits.Knight, card.Traits.Spirit),
	card.WithAbility(
		card.Trigger.Play, card.Stun{Target: card.Target.Creature}),
	card.WithAbility(
		card.Trigger.Action, card.PutFromPlay{
			Target:      card.Target.This,
			Destination: card.To.Hand,
		}),
)
