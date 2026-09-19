package massmutation

import (
	"github.com/dmikalova/vactrol/internal/card"
	"github.com/dmikalova/vactrol/internal/cards/clusters"
)

// It's Coming...
//
//	House:  None
//	Type:   Tactic
//	Rarity: Special
//	Bonus:  Æmber
//
//	Play: Search your deck and discard pile for either half of a gigantic creature, reveal it, and put it into your hand. Shuffle your deck.
var ItsComing = set.New(
	"It's Coming...",
	card.House.None,
	card.Type.Tactic,
	card.Rarity.Special,
	card.Provenance(card.MM, "117"),
	card.InCluster(clusters.Tutors),
	card.Houseless(),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{
			Effects: []card.Effect{
				card.Search{
					Sources: []card.Zone{card.Deck, card.Discard},
					Filter:  card.Filter{Gigantic: true},
					Reveal:  true,
					Dest:    card.To.Hand,
				},
				card.Shuffle{},
			},
		}),
)
