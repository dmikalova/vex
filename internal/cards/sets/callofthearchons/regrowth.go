package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Regrowth
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Put a creature from your discard pile into your hand.
var Regrowth = set.New(
	"Regrowth",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.CotA, "332"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.PutCard{Zones: []card.Zone{card.Discard},
			Selection:   card.Chosen{Type: card.Type.Creature},
			Destination: card.To.Hand,
		}),
)
