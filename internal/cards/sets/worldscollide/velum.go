package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Velum
//
//	House:  Logos
//	Type:   Creature
//	Rarity: Special
//	Power:  2
//	Traits: Human • Scientist
//
//	Reap: Archive a card from your hand, or 2 cards if you control Hyde.
//	Destroyed: Archive Hyde from your discard pile -> archive Velum from play.
var Velum = card.New(
	"Velum",
	card.House.Logos,
	card.Type.Creature,
	card.Rarity.Special,
	card.Provenance(card.WC, "181"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Human, card.Traits.Scientist),
	card.WithAbility(card.Trigger.Reap, card.ArchiveFromHand{
		Amount: 1,
		Or: card.OrAmount{
			Amount: 2,
			When:   card.ControlsNamed{Name: "Hyde"},
		},
	}),
	card.WithAbility(card.Trigger.Destroyed, card.Then{
		First:  card.ArchiveFromDiscard{Name: "Hyde"},
		Result: card.ArchiveFromPlay{Target: card.Target.This},
	}),
)
