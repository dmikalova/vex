package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Spare Arm Carmine
//
//	House:  Shadows
//	Type:   Creature
//	Rarity: Common
//	Power:  2
//	Traits: Mutant
//
//	Reap: Steal 1 Æmber. If you control more Mutant creatures than your opponent, steal 1 Æmber.
var SpareArmCarmine = set.New(
	"Spare Arm Carmine",
	card.House.Shadows,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MoMu, "307"),
	card.WithPower(2),
	card.WithTraits(card.Traits.Mutant),
	card.WithAbility(
		card.Trigger.Reap, card.Sequence{Effects: []card.Effect{
			card.StealAember{Amount: 1},
			card.Conditional{
				Cond: card.ControlsMoreCreatures{Trait: card.Traits.Mutant},
				Then: card.StealAember{Amount: 1},
			},
		}}),
)
