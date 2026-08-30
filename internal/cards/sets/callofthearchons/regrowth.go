package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Regrowth
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Put a creature from your discard pile into your hand.
var Regrowth = card.New(
	"Regrowth",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, 332),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.PutFromDiscard{
			Type:        card.Type.Creature,
			Destination: card.To.Hand,
		}),
)
