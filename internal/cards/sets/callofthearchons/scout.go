package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Scout
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//	Æmber:  1
//
//	Play: Give skirmish to, ready, and fight with up to 2 different friendly creatures, one at a time.
var Scout = card.New(
	"Scout",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, 334),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.OneAtATime{
			Times:  2,
			Target: card.Target.FriendlyCreature,
			Verbs: []card.CreatureVerb{
				card.GainKeywordVerb{Keyword: card.Keyword.Skirmish},
				card.ReadyVerb{},
				card.FightVerb{},
			},
		}),
)
