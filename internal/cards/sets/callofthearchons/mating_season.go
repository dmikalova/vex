package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Mating Season
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Shuffle each Mars creature into its owner's deck. For each creature shuffled into their deck this way, each player gains 1 Æmber.
var MatingSeason = set.New(
	"Mating Season",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "170"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{Effects: []card.Effect{
			card.PutFromPlay{
				Target:      card.Target.EachCreature.House(card.Houses.Named(card.House.Self)),
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
