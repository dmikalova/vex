package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Z.Y.X. Researcher
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Human • Scientist
//
//	Play: Choose one:
//	- Archive the top card of your deck
//	- Archive the top card of your discard pile.
var ZYXResearcher = set.New(
	"Z.Y.X. Researcher",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.AoA, "123"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Human, card.Traits.Scientist),
	card.WithAbility(
		card.Trigger.Play, card.ChooseOne{
			Options: []card.Effect{
				card.ArchiveCard{Zone: card.Deck, Selection: card.Top{}},
				card.ArchiveCard{Zone: card.Discard, Selection: card.Top{}},
			},
		}),
)
