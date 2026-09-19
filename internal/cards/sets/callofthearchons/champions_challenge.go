package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Champion's Challenge
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Destroy each enemy creature except the most powerful enemy creature. Destroy each friendly creature except the most powerful friendly creature. Ready and fight with a friendly creature.
var ChampionsChallenge = set.New(
	"Champion's Challenge",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "6"),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{Effects: []card.Effect{
			card.Destroy{
				Target: card.Target.EachEnemyCreature.Refine(card.Except(card.MostPowerful)),
			},
			card.Destroy{
				Target: card.Target.EachFriendlyCreature.Refine(card.Except(card.MostPowerful)),
			},
			card.OnChooseCreature{
				Target: card.Target.FriendlyCreature,
				Verbs:  []card.CreatureVerb{card.ReadyVerb{}, card.FightVerb{}},
			},
		}}),
)
