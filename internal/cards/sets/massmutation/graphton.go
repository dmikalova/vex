package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Graphton
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Robot • Mutant
//
//	Reap: Archive the top card of your deck.
var Graphton = set.New(
	"Graphton",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MoMu, "115"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Robot, card.Traits.Mutant),
	card.WithAbility(
		card.Trigger.Reap, card.ArchiveCard{
			Zone:      card.Deck,
			Selection: card.Top{},
			Amount:    1,
		}),
)
