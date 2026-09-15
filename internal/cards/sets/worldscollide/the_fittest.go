package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// The Fittest
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Give each friendly creature a +1 power counter.
var TheFittest = set.New(
	"The Fittest",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "366"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.AddPowerCounter{
			Target: card.Target.EachFriendlyCreature,
			Amount: 1,
		}),
)
