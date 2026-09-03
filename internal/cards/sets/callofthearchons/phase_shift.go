package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Phase Shift
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Play a non-Logos card.
var PhaseShift = card.New(
	"Phase Shift",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 117),
	card.WithAbility(
		card.Trigger.Play, card.PlayFrom{
			From:   card.Hand,
			House:  card.House.Self,
			Except: true,
		}),
)
