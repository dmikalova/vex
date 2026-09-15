package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Lost in the Woods
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Shuffle 2 friendly creatures into their owners' decks, and shuffle 2 enemy creatures into their owners' decks.
var LostInTheWoods = set.New(
	"Lost in the Woods",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, "327"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.PutChosen{
				Amount:      2,
				Target:      card.Target.EachFriendlyCreature,
				Destination: card.To.DeckShuffled,
			},
			card.PutChosen{
				Amount:      2,
				Target:      card.Target.EachEnemyCreature,
				Destination: card.To.DeckShuffled,
			},
		}}),
)
