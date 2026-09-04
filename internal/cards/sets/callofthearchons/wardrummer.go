package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Wardrummer
//
//	House:  Brobnar
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Goblin
//
//	Play: Put each other friendly Brobnar creature into its owner's hand.
var Wardrummer = card.New(
	"Wardrummer",
	card.House.Brobnar,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 49),
	card.WithPower(3),
	card.WithTraits(card.Traits.Goblin),
	card.WithAbility(
		card.Trigger.Play, card.PutFromPlay{
			Target:      card.Target.EachOtherFriendlyCreature.OfHouse(card.House.Brobnar),
			Destination: card.To.Hand,
		}),
)
