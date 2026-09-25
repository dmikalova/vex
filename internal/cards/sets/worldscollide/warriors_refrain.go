package worldscollide

import "github.com/dmikalova/vex/internal/card"

// Warriors' Refrain
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Stun each creature with power 3 or lower.
var WarriorsRefrain = set.New(
	"Warriors' Refrain",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "16"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Stun{
			Target: card.Target.EachCreature.PowerAtMost(3),
		}),
)
