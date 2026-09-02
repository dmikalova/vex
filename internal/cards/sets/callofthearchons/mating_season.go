package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Mating Season
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Shuffle each Mars creature into its owner's deck. For each creature shuffled into their deck this way, each player gains 1 Æmber.
var MatingSeason = card.New(
	"Mating Season",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 170),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.Sentence{Effect: card.PutFromPlay{
				Target:      card.Target.EachCreature.OfHouse(card.House.Mars),
				Destination: card.To.DeckShuffled,
			}},
			card.Sentence{Effect: card.GainAember{
				Player: card.EachPlayer,
				Amount: 1,
				Per:    card.CreaturesShuffledIntoDeckThisWay{Player: card.Controller},
			}},
		}}),
)
