package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Tormax
//
//	House:  Dis
//	Type:   Creature
//	Rarity: Rare
//	Power:  8
//	Traits: Demon
//
//	Play/Fight/Reap: Discard your hand, and purge 2 random cards from your opponent's hand.
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
				Zone:      card.Hand,
				Player:    card.Opponent,
				Selection: card.Random{Count: 2},
			},
		}}),
)
