package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Soft Landing
//
//	House:  Mars
//	Type:   Tactic
//	Rarity: Common
//
//	Play: The next creature or artifact you play this turn enters play ready.
var SoftLanding = card.New(
	"Soft Landing",
	card.House.Mars,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 177),
	card.WithAbility(
		card.Trigger.Play, card.NextPlayed{
			Type:       card.Type.Any,
			EntersPlay: card.Ready{Target: card.Target.Triggering},
		}),
)
