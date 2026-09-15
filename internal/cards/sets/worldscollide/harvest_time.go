package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Harvest Time
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Choose a creature - purge each creature that shares a trait with it. For each card they controlled that was purged this way, each player gains 1 Æmber.
var HarvestTime = set.New(
	"Harvest Time",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.WC, "106"),
	card.WithAbility(
		card.Trigger.Play, card.ChooseCreatureThen{
			Target: card.Target.Creature,
			Then: card.Sentences{Effects: []card.Effect{
				card.PurgeCreature{Target: card.Target.EachCreature.SharingTrait()},
				card.GainAember{
					Player: card.EachPlayer,
					Amount: 1,
					Per: card.ProducedThisWay{
						Tally:  card.Tally.CardsPurged,
						Player: card.Controller,
					},
				},
			}},
		}),
)
