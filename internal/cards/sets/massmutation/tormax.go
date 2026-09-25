package massmutation

import "github.com/dmikalova/vex/internal/card"

// Tormax
//
//	House:  Dis
//	Type:   Gigantic Creature
//	Rarity: Rare
//	Power:  8
//	Traits: Demon
//
//	Play/Fight/Reap: Discard your hand. Your opponent purges 2 random cards from their hand.
var Tormax = set.Gigantic(
	"Tormax",
	card.House.Dis,
	card.Rarity.Rare,
	card.Provenance(card.MoMu, "020"),
	card.WithPower(8),
	card.WithTraits(card.Traits.Demon),
	card.WithAbility(
		card.Trigger.PlayFightReap, card.Sequence{Effects: []card.Effect{
			card.DiscardHand{Player: card.Controller},
			card.PurgeCard{
				Zones:     []card.Zone{card.Hand},
				Player:    card.Opponent,
				Selection: card.Random{},
				Quantity:  card.Takes{N: card.Fixed(2)},
			},
		}}),
)
