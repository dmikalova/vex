package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Hecatomb
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Destroy each Dis creature. For each creature they controlled that was destroyed this way, each player gains 1 Æmber.
var Hecatomb = set.New(
	"Hecatomb",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "63"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.Destroy{
				Target: card.Target.EachCreature.House(card.Houses.Named(card.House.Self)),
			},
			card.GainAember{
				Player: card.EachPlayer,
				Amount: 1,
				Per: card.ProducedThisWay{
					Tally:  card.Tally.CreaturesDestroyed,
					Player: card.Controller,
				},
			},
		}}),
)
