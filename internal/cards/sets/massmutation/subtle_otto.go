package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Subtle Otto
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  1
//	Traits: Mutant • Thief
//
//	Play: Your opponent discards a random card from their hand.
var SubtleOtto = set.New(
	"Subtle Otto",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "258"),
	card.WithPower(1),
	card.WithTraits(card.Traits.Mutant, card.Traits.Thief),
	card.WithAbility(
		card.Trigger.Play, card.DiscardCard{
			Player:    card.Opponent,
			Zone:      card.Hand,
			Selection: card.Random{},
		}),
)
