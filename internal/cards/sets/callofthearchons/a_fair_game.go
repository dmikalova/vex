package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// A Fair Game
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Discard the top card of your opponent's deck. Reveal your opponent's hand. For each card of the discarded card's house revealed this way, gain 1 Æmber. Discard the top card of your deck. Reveal your hand. For each card of the discarded card's house revealed this way, your opponent gains 1 Æmber.
var AFairGame = card.New(
	"A Fair Game",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 53),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.Sentence{Effect: card.DiscardTopOfDeck{Player: card.Opponent}},
			card.Sentence{Effect: card.RevealHand{Player: card.Opponent}},
			card.Sentence{Effect: card.GainAember{
				Player: card.Controller,
				Amount: 1,
				Per:    card.CardsInHand{Player: card.Opponent, House: card.TheContextualHouse},
			}},
			card.Sentence{Effect: card.DiscardTopOfDeck{Player: card.Controller}},
			card.Sentence{Effect: card.RevealHand{Player: card.Controller}},
			card.Sentence{Effect: card.GainAember{
				Player: card.Opponent,
				Amount: 1,
				Per:    card.CardsInHand{Player: card.Controller, House: card.TheContextualHouse},
			}},
		}}),
)
