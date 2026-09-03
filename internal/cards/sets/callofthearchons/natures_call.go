package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Nature's Call
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Put up to 3 creatures into their owners' hands.
var NaturesCall = card.New(
	"Nature's Call",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 329),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.PutChosen{
			Count:       3,
			UpTo:        true,
			Target:      card.Target.EachCreature,
			Destination: card.To.Hand,
		}),
)
