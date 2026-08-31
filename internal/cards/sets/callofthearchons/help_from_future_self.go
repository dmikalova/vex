package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Help from Future Self
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: FIXED
//	Æmber:  1
//
//	Play: Search your deck and discard pile for a Timetraveller, reveal it, and put it into your hand, and shuffle your discard pile into your deck.
var HelpFromFutureSelf = card.New(
	"Help from Future Self",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Fixed,
	card.Provenance(card.CotA, 111),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{
			Effects: []card.Effect{
				card.SearchForName{Name: "Timetraveller"},
				card.ShuffleIntoDeck{Zones: []card.Zone{card.Discard}},
			},
		}),
)
