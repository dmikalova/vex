package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Scout
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Rare
//	Bonus:  Æmber
//
//	Play: Give skirmish to, ready, and fight with up to 2 different friendly creatures, one at a time.
var Scout = set.New(
	"Scout",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.CotA, "334"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.OneAtATime{
			Times:  card.Fixed(2),
			Target: card.Target.FriendlyCreature,
			Verbs: []card.CreatureVerb{
				card.GainKeywordVerb{Keyword: card.Keyword.Skirmish},
				card.ReadyVerb{},
				card.FightVerb{},
			},
		}),
)
