package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Random Access Archives
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Archive the top card of your deck.
var RandomAccessArchives = set.New(
	"Random Access Archives",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "119"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play,
		card.ArchiveCard{
			Zone:      card.Deck,
			Selection: card.Top{},
		}),
)
