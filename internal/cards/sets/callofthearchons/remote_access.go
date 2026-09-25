package callofthearchons

import "github.com/dmikalova/vex/internal/card"

// Remote Access
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Use an enemy artifact.
var RemoteAccess = set.New(
	"Remote Access",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "120"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Use{
			Max:    1,
			Target: card.Target.EachEnemyArtifact,
		}),
)
