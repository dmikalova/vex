package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Master the Theory
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: If there are no friendly creatures in play, for each enemy creature in play, you may archive a card from your hand.
var MasterTheTheory = card.New(
	"Master the Theory",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "148"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Conditional{
			Cond: card.InPlay{
				Player: card.Controller,
				Type:   card.Type.Creature,
				None:   true,
			},
			Then: card.ForEach{
				Times: card.InPlay{
					Player: card.Opponent,
					Type:   card.Type.Creature,
				},
				Do: card.May{Do: card.ArchiveFromHand{Amount: 1}},
			},
		}),
)
