package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Library Access
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//
//	Play: For the remainder of the turn, each time you play another card, draw a card. Purge Library Access.
var LibraryAccess = card.New(
	"Library Access",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 115),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{Effects: []card.Effect{
			card.ForRemainderOfTurn{
				On: card.Event.CardPlayed,
				Do: card.Draw{Amount: 1},
			},
			card.PurgeSource{},
		}}),
)
