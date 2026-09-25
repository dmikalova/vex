package massmutation

import "github.com/dmikalova/vex/internal/card"

// Wail of the Damned
//
//	House:  Dis
//	Type:   Tactic
//	Rarity: Uncommon
//
//	Play: Destroy a creature with no bonus icons.
//	Enhance Capture.
var WailOfTheDamned = set.New(
	"Wail of the Damned",
	card.House.Dis,
	card.Type.Tactic,
	card.Rarity.Uncommon,
	card.Provenance(card.MM, "033"),
	card.WithEnhance(card.Bonus.Capture),
	card.WithAbility(
		card.Trigger.Play, card.Destroy{
			Target: card.Target.Creature.WithoutBonusIcons(),
		}),
)
