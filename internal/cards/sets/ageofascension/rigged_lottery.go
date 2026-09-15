package ageofascension

import "github.com/dmikalova/vactrol/internal/card"

// Rigged Lottery
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Discard the top 5 cards of each player's deck. For each Shadows card discarded this way, its owner gains 1 Æmber.
var RiggedLottery = set.New(
	"Rigged Lottery",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.AoA, "309"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{Effects: []card.Effect{
			card.DiscardTop{Player: card.EachPlayer, Amount: 5},
			card.ForEachDiscarded{
				House: card.House.Self,
				Do:    card.GainAember{Player: card.ItsOwner, Amount: 1},
			},
		}}),
)
