package ageofascension

import "github.com/dmikalova/vex/internal/card"

// Director of Z.Y.X.
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  3
//	Traits: Cyborg • Scientist
//
//	Elusive.
//	At the start of your turn, archive the top card of your deck.
var DirectorOfZYX = set.New(
	"Director of Z.Y.X.",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.AoA, "127"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Cyborg, card.Traits.Scientist),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.StartOfTurn,
		card.ArchiveCard{
			Zone:      card.Deck,
			Selection: card.Top{},
		}),
)
