package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Sloppy Labwork
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Archive a card from your hand. Discard a card from your hand.
var SloppyLabwork = set.New(
	"Sloppy Labwork",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "123"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(card.Trigger.Play, card.Sequence{
		Effects: []card.Effect{
			card.ArchiveCard{
				Zone:      card.Hand,
				Selection: card.Chosen{},
			},
			card.DiscardCard{
				Player:    card.Controller,
				Zones:     []card.Zone{card.Hand},
				Selection: card.Chosen{},
			},
		},
	}),
)
