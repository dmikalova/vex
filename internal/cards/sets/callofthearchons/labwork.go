package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Labwork
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Archive a card from your hand.
var Labwork = set.New(
	"Labwork",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, "114"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play,
		card.ArchiveCard{
			Zone:      card.Hand,
			Selection: card.Chosen{},
		},
	),
)
