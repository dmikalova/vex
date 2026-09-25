package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Research Smoko
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Uncommon
//	Power:  2
//	Traits: Mutant
//
//	Destroyed: Archive the top card of your deck.
var ResearchSmoko = set.New(
	"Research Smoko",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "151"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Mutant),
	card.WithAbility(
		card.Trigger.Destroyed,
		card.ArchiveCard{
			Zone:      card.Deck,
			Selection: card.Top{},
		}),
)
