package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Daemo-Saurus
//
//	House:  Saurian
//	Type:   Creature
//	Rarity: Common
//	Power:  4
//	Traits: Mutant • Dinosaur
//
//	Play: You may exalt Daemo-Saurus -> deal 3 damage to a creature.
//	Destroyed: Steal 1 Æmber.
var DaemoSaurus = set.New(
	"Daemo-Saurus",
	card.House.Saurian,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "190"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Mutant, card.Traits.Dinosaur),
	card.WithAbility(
		card.Trigger.Play, card.May{
			Do: card.Then{
				First: card.Exalt{Target: card.Target.This, Amount: 1},
				Result: card.DealDamage{
					Amount: 3,
					Target: card.Target.Creature,
				},
			},
		}),
	card.WithAbility(
		card.Trigger.Destroyed, card.StealAember{Amount: 1}),
)
