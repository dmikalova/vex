package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Battle Fleet
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Reveal any number of Mars cards from your hand, and for each card revealed this way, draw a card.
var BattleFleet = card.New(
	"Battle Fleet",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 161),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.Reveal{
				Player: card.Controller,
				House:  card.House.Mars,
			},
			card.Draw{
				Amount: 1,
				Per:    card.CardsRevealed{},
			},
		}}),
)
