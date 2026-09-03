package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Lost in the Woods
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Shuffle 2 friendly creatures into their owners' decks, and shuffle 2 enemy creatures into their owners' decks.
var LostInTheWoods = card.New(
	"Lost in the Woods",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 327),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.PutChosen{
				Count:       2,
				Target:      card.Target.EachFriendlyCreature,
				Destination: card.To.DeckShuffled,
			},
			card.PutChosen{
				Count:       2,
				Target:      card.Target.EachEnemyCreature,
				Destination: card.To.DeckShuffled,
			},
		}}),
)
