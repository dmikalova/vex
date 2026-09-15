package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Radiant Truth
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Uncommon
//	Bonus:  Æmber
//
//	Play: Stun each enemy creature that is not on a flank.
var RadiantTruth = set.New(
	"Radiant Truth",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.CotA, "224"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play,
		card.Stun{Target: card.Target.EachEnemyCreature.NotOnFlank()},
	),
)
