package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Battle Fleet
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Reveal any number of Mars cards from your hand, and for each card revealed this way, draw a card.
var BattleFleet = set.New(
	"Battle Fleet",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "161"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.RevealHand{
				Player: card.Controller,
				House:  card.Houses.Named(card.House.Self),
			},
			card.Draw{
				Amount: 1,
				Per:    card.CardsRevealed{},
			},
		}}),
)
