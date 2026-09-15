package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Mind Over Matter
//
//	House:  Logos
//	Type:   Tactic
//	Rarity: Rare
//
//	Play: Archive each creature from play.
var MindOverMatter = set.New(
	"Mind Over Matter",
	card.House.Logos,
	card.Type.Tactic,
	card.Rarity.Rare,
	card.Provenance(card.MM, "107"),
	card.WithAbility(
		card.Trigger.Play, card.ArchiveFromPlay{
			Target: card.Target.EachCreature,
		}),
)
