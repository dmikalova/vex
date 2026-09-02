package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Remote Access
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//	Æmber:  1
//
//	Play: Use an enemy artifact.
var RemoteAccess = card.New(
	"Remote Access",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, 120),
	card.WithAemberBonus(1),
	card.WithAbility(
		card.Trigger.Play, card.Use{
			Max:    1,
			Target: card.Target.EachEnemyArtifact,
		}),
)
