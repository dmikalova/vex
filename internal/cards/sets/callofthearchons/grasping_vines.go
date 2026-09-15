package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Grasping Vines
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Put up to 3 artifacts into their owners' hands.
var GraspingVines = set.New(
	"Grasping Vines",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "324"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play,
		card.PutChosen{
			Amount:      3,
			UpTo:        true,
			Target:      card.Target.EachArtifact,
			Destination: card.To.Hand,
		},
	),
)
