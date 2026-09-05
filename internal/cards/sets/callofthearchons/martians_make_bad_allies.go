package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Martians Make Bad Allies
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Reveal your hand, and purge each non-Mars creature from your hand, and for each creature purged this way, gain 1 Æmber.
var MartiansMakeBadAllies = card.New(
	"Martians Make Bad Allies",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 168),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.RevealHand{Player: card.Controller},
			card.PurgeEachFromHand{
				Player:      card.Controller,
				Type:        card.Type.Creature,
				ExceptHouse: card.House.Self,
			},
			card.GainAember{
				Player: card.Controller,
				Amount: 1,
				Per:    card.CardsPurged{},
			},
		}}),
)
