package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Reclaimed by Nature
//
//	House:  Untamed
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Purge an artifact. Resolve that card's bonus icons.
var ReclaimedByNature = set.New(
	"Reclaimed by Nature",
	card.House.Untamed,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "374"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Sentences{Effects: []card.Effect{
			card.PurgeCreature{Target: card.Target.Artifact},
			card.ResolveBonusIcons{Target: card.Target.Triggering},
		}}),
)
