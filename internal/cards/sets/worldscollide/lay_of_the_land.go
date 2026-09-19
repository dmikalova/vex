package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Lay of the Land
//
//	House:  Star Alliance
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Look at the top 3 cards of your deck and put them back in any order. Draw a card.
var LayOfTheLand = set.New(
	"Lay of the Land",
	card.House.StarAlliance,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.WC, "313"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(card.Trigger.Play, card.Sequence{Effects: []card.Effect{
		card.LookAtTopOfDeck{
			Amount: 3,
			Then:   []card.TopAct{card.ReorderRest{}},
		},
		card.Draw{Amount: 1},
	}}),
)
