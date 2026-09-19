package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Lesser Oxtet
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  3
//	Traits: Demon
//
//	Elusive.
//	Play: Purge each card from your hand.
//	Reap: Keys cost +3 Æmber during your opponent's next turn.
var LesserOxtet = set.New(
	"Lesser Oxtet",
	card.House.Dis,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.WC, "109"),
	card.WithPower(3),
	card.WithTraits(card.Traits.Demon),
	card.WithKeywords(card.Keyword.Elusive),
	card.WithAbility(
		card.Trigger.Play,
		card.PurgeCard{
			Zones:     []card.Zone{card.Hand},
			Player:    card.Controller,
			Selection: card.Each{},
		},
	),
	card.WithAbility(
		card.Trigger.Reap, card.RaiseKeyCost{
			Player:   card.Opponent,
			Amount:   3,
			Duration: card.Duration.OpponentNextTurn,
		}),
)
