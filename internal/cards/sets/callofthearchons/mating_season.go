package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Mating Season
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Shuffle each Mars Creature into its owner's deck. For each Creature shuffled into their deck this way, each player gains 1 Æmber.
var MatingSeason = card.New(
	"Mating Season",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "170"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{Effects: []card.Effect{
			card.PutFromPlay{
				Target:      card.Target.EachCreature.OfHouse(card.House.Self),
				Destination: card.To.DeckShuffled,
			},
			card.GainAember{
				Player: card.EachPlayer,
				Amount: 1,
				Per: card.ProducedThisWay{
					Tally:  card.Tally.CreaturesShuffledIntoDeck,
					Player: card.Controller,
				},
			},
		}}),
)
