package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// One Last Job
//
//	House:  Shadows
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Purge each friendly Shadows creature. For each creature purged this way, steal 1 Æmber.
var OneLastJob = card.New(
	"One Last Job",
	card.House.Shadows,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 277),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{Effects: []card.Effect{
			card.PurgeCreature{
				Target: card.Target.EachFriendlyCreature.OfHouse(card.House.Shadows),
			},
			card.StealAember{
				Amount: 1,
				Per:    card.CardsPurged{},
			},
		}}),
)
