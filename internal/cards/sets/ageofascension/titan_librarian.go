package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Titan Librarian
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Cyborg • Scientist
//
//	At the end of your turn, if Titan Librarian is not on a flank, archive a card from your hand.
var TitanLibrarian = set.New(
	"Titan Librarian",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, "120"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Cyborg, card.Traits.Scientist),
	card.WithAbility(
		card.Trigger.EndOfTurn, card.Conditional{
			Cond: card.Not{Cond: card.OnFlank{}},
			Then: card.ArchiveCard{
				Zone:      card.Hand,
				Selection: card.Chosen{},
			},
		}),
)
