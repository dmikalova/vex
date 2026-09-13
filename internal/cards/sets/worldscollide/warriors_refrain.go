package worldscollide

import "github.com/dmikalova/vactrol/internal/card"

// Warriors' Refrain
//
//	House:  Brobnar
//	Type:   Tactic
//	Rarity: Common
//	Æmber:  1
//
//	Play: Stun each Creature with power 3 or lower.
var WarriorsRefrain = card.New(
	"Warriors' Refrain",
	card.House.Brobnar,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.WC, "16"),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Stun{
			Target: card.Target.EachCreature.PowerAtMost(3),
		}),
)
