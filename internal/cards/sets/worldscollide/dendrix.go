package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Dendrix
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Common
//	Power:  5
//	Traits: Demon
//
//	Fight: Your opponent discards a random card from their hand.
var Dendrix = set.New(
	"Dendrix",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.WC, "71"),
	card.WithPower(5),
	card.WithTraits(card.Traits.Demon),
	card.WithAbility(
		card.Trigger.Fight,
		card.DiscardCard{
			Player:    card.Opponent,
			Zones:     []card.Zone{card.Hand},
			Selection: card.Random{},
		},
	),
)
