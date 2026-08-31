package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Champion's Challenge
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Destroy each enemy creature except the most powerful enemy creature and each friendly creature except the most powerful friendly creature, and ready and fight with a friendly creature.
var ChampionsChallenge = card.New(
	"Champion's Challenge",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 6),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.Destroy{Target: card.Target.EachEnemyCreature.Selector(card.ExceptMostPowerful)},
			card.Destroy{
				Target: card.Target.EachFriendlyCreature.Selector(card.ExceptMostPowerful),
			},
			card.OnChooseCreature{
				Target: card.Target.FriendlyCreature,
				Verbs:  []card.CreatureVerb{card.ReadyVerb{}, card.FightVerb{}},
			},
		}}),
)
