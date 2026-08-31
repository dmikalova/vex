package callofthearchons

import "github.com/dmikalova/vactrol/internal/card"

// Vezyma Thinkdrone
//
//	House:  Mars
//	Type:   Creature
//	Rarity: Common
//	Power:  3
//	Traits: Martian • Scientist
//
//	Reap: You may archive a friendly creature or artifact from play.
var VezymaThinkdrone = card.New(
	"Vezyma Thinkdrone",
	card.House.Mars,
	card.Type.Creature,
	card.Rarity.Common,
	card.Provenance(card.CotA, 202),
	card.WithPower(3),
	card.WithTraits("Martian", "Scientist"),
	card.WithAbility(
		card.Trigger.Reap, card.May{
			Do: card.ArchiveFromPlay{Target: card.Target.FriendlyInPlay},
		}),
)
