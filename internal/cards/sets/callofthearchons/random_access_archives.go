package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Random Access Archives
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Archive the top card of your deck.
var RandomAccessArchives = set.New(
	"Random Access Archives",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "119"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play,
		card.ArchiveCard{Zone: card.Deck, Selection: card.Top{}}),
)
