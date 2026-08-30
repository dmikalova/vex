package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// A Fair Game
//
//	House:  Dis
//	Type:   Action
//	Rarity: Rare
//
//	Play: Discard the top card of your opponent's deck and reveal their hand. You gain 1 Æmber for each card of the discarded card's house revealed this way. Discard the top card of your deck and reveal your hand. Your opponent gains 1 Æmber for each card of the discarded card's house revealed this way.
var AFairGame = card.New(
	"A Fair Game",
	card.House.Dis,
	card.Type.Action,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 53),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.DiscardTopOfDeckAndRevealHandForAember{
				Player: card.Opponent,
				Gainer: card.Controller,
			},
			card.DiscardTopOfDeckAndRevealHandForAember{
				Player: card.Controller,
				Gainer: card.Opponent,
			},
		}}),
)
