package massmutation

import "github.com/dmikalova/vactrol/internal/card"

// Commander Dhrxgar
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Rare
//	Power:  4
//	Traits: Mutant
//
//	After an upgrade enters play, if it is attached to Commander Dhrxgar or one of its neighbors, gain 1 Æmber.
var CommanderDhrxgar = set.New(
	"Commander Dhrxgar",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Rare,
	card.Provenance(card.MM, "337"),
	card.WithPower(4),
	card.WithTraits(card.Traits.Mutant),
	card.WithAbility(
		card.Trigger.AfterUpgradeEnters, card.Conditional{
			Cond: card.ItAttachedToThisOrNeighbor{},
			Then: card.GainAember{
				Player: card.Controller,
				Amount: 1,
			},
		}),
)
