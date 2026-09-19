package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Mindfire
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Your opponent discards a random card from their hand. For each bonus icon on the discarded card, steal 1 Æmber.
var Mindfire = set.New(
	"Mindfire",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "012"),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{Effects: []card.Effect{
			card.DiscardCard{
				Player:    card.Opponent,
				Zones:     []card.Zone{card.Hand},
				Selection: card.Random{},
				Bind:      true,
			},
			card.StealAember{
				Amount: 1,
				Per:    card.BonusIconsOfChosen{Noun: card.ItNoun.DiscardedCard},
			},
		}}),
)
