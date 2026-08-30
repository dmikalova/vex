package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Hysteria
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Put each creature into its owner's hand.
var Hysteria = card.New(
	"Hysteria",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 65),
	card.WithAbility(
		card.Trigger.Play, card.PutFromPlay{
			Target:      card.Target.EachCreature,
			Destination: card.To.Hand,
		}),
)
