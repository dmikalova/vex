package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Tautau Vapors
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Draw 2 cards. Archive a card from your hand.
var TautauVapors = set.New(
	"Tautau Vapors",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "139"),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{
			Effects: []card.Effect{
				card.Draw{Amount: 2},
				card.ArchiveCard{
					Zone:      card.Hand,
					Selection: card.Chosen{},
				},
			},
		}),
)
