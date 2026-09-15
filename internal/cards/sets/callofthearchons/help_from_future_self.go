package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Help from Future Self
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Connected
//	Bonus:  Æmber
//
//	Play: Search your deck and discard pile for a Timetraveller, reveal it, and put it into your hand, and shuffle your discard pile into your deck.
var HelpFromFutureSelf = set.New(
	"Help from Future Self",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Connected,
	card.Provenance(card.CotA, "111"),
	card.InCluster(timetravellerCluster),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{
			Effects: []card.Effect{
				card.SearchForName{Name: Timetraveller.Name},
				card.Shuffle{Zones: []card.Zone{card.Discard}},
			},
		}),
)
