package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Curse of Vanity
//
//	House:  Saurian
//	Type:   Tactic
//	Rarity: Common
//	Bonus:  Æmber
//
//	Play: Exalt a friendly creature and an enemy creature.
var CurseOfVanity = set.New(
	"Curse of Vanity",
	card.House.Saurian,
	card.Type.Tactic,
	card.Rarity.Common,
	card.Provenance(card.MM, "189"),
	card.WithBonus(card.Bonus.Aember),
	card.WithAbility(
		card.Trigger.Play, card.Sequence{Effects: []card.Effect{
			card.Exalt{Target: card.Target.FriendlyCreature, Amount: 1},
			card.Exalt{Target: card.Target.EnemyCreature, Amount: 1},
		}}),
)
