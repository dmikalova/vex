package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Fangs of Gizelhart
//
//	House:  Sanctum
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Purge the most powerful creature.
var FangsOfGizelhart = set.New(
	"Fangs of Gizelhart",
	card.House.Sanctum,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "133"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.PurgeCreature{
			Target: card.Target.EachCreature.Refine(card.MostPowerful),
		}),
)
