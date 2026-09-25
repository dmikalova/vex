package massmutation

import "github.com/dmikalova/vex/internal/card"

// Crewman Jorg
//
//	House:  Star Alliance
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Human • Thief
//
//	Action: If Crewman Jorg has no Star Alliance neighbor, steal 1 Æmber.
//	Enhance Capture.
var CrewmanJorg = set.New(
	"Crewman Jorg",
	card.House.StarAlliance,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.MM, "305"),
	card.WithEnhance(card.Bonus.Capture),
	card.WithPower(3),
	card.WithTraits(card.Traits.Human, card.Traits.Thief),
	card.WithAbility(
		card.Trigger.Action, card.Conditional{
			Cond: card.SourceHasNoNeighbor{House: card.Houses.Named(card.House.Self)},
			Then: card.StealAember{Amount: 1},
		}),
)
