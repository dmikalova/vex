package massmutation

import "github.com/dmikalova/vex/internal/card"

// Eclectic Inquiry
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Archive the top 2 cards of your deck.
var EclecticInquiry = set.New(
	"Eclectic Inquiry",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "071"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.ArchiveCard{
			Zone:      card.Deck,
			Selection: card.Top{},
			Quantity:  card.Takes{N: card.Fixed(2)},
		}),
)
