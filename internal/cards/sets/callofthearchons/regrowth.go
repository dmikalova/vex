package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Regrowth
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Put a Creature from your discard pile into your hand.
var Regrowth = set.New(
	"Regrowth",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, "332"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.PutFromDiscard{
			Selection:   card.Chosen{Type: card.Type.Creature},
			Destination: card.To.Hand,
		}),
)
