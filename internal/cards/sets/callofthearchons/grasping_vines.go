package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Grasping Vines
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Put up to 3 artifacts into their owners' hands.
var GraspingVines = card.New(
	"Grasping Vines",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 324),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play,
		card.PutChosen{
			Count:       3,
			UpTo:        true,
			Target:      card.Target.EachArtifact,
			Destination: card.To.Hand,
		},
	),
)
