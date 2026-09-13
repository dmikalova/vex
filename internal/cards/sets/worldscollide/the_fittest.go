package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// The Fittest
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Give each friendly Creature a +1 power counter.
var TheFittest = card.New(
	"The Fittest",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "366"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.AddPowerCounter{
			Target: card.Target.EachFriendlyCreature,
			Amount: 1,
		}),
)
