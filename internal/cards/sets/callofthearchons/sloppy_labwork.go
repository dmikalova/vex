package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Sloppy Labwork
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Archive a card from your hand. Discard a card from your hand.
var SloppyLabwork = card.New(
	"Sloppy Labwork",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 123),
	card.WithAemberBonus(1),
	card.WithAbility(card.Trigger.Play, card.Sequence{
		Effects: []card.Effect{
			card.Sentence{Effect: card.ArchiveFromHand{Count: 1}},
			card.Sentence{Effect: card.DiscardFromHand{Count: 1}},
		},
	}),
)
