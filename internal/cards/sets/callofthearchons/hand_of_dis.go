package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Hand of Dis
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Common
//
//	Play: Destroy a creature that is not on a flank.
var HandOfDis = card.New(
	"Hand of Dis",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 62),
	card.WithAbility(card.Trigger.Play, card.Destroy{Target: card.Target.Creature.NotOnFlank()}),
)
