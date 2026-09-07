package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Lay of the Land
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Look at the top 3 cards of your deck and put them back in any order, and draw a card.
var LayOfTheLand = card.New(
	"Lay of the Land",
	card.House.StarAlliance,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, 313),
	card.WithAemberBonus(1),
	card.WithAbility(card.Trigger.Play, card.Sequence{Effects: []card.Effect{
		card.ReorderTop{Amount: 3},
		card.Draw{Amount: 1},
	}}),
)
